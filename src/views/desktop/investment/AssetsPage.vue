<template>
    <div class="page-content">
        <div class="page-body">
            <AssetSearchBar
                :show-manage-button="showManageButton"
                :watchlist-asset-ids="watchlistAssetIds"
                @manage="adminDialog = true"
                @select="onSearchResultSelect"
                @watchlist-updated="reloadWatchlist"
            />

            <v-card>
                <v-tabs v-model="activeTab" class="px-4">
                    <v-tab value="holdings">
                        {{ tt('Holdings') }}
                        <v-chip size="small" variant="tonal" class="ml-2">{{ holdingsCount }}</v-chip>
                    </v-tab>
                    <v-tab value="watchlist">
                        {{ tt('Watchlist') }}
                        <v-chip size="small" variant="tonal" class="ml-2">{{ watchlist.length }}</v-chip>
                    </v-tab>
                </v-tabs>

                <HoldingsTab v-if="activeTab === 'holdings'" :loading="loading" />

                <WatchlistTab
                    v-if="activeTab === 'watchlist'"
                    :loading="loading"
                    :watchlist="watchlist"
                    @select="onWatchlistSelect"
                    @removed="reloadWatchlist"
                />

                <template v-if="activeTab === 'holdings'">
                    <v-divider />
                    <HoldingsSummaryBar />
                </template>
            </v-card>

            <AssetDetailDialog
                v-model="detailDialog"
                :item="selectedItem"
            />

            <AssetAdminDialog v-model="adminDialog" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';

import AssetSearchBar from './assets/AssetSearchBar.vue';
import HoldingsTab from './assets/HoldingsTab.vue';
import WatchlistTab from './assets/WatchlistTab.vue';
import HoldingsSummaryBar from './assets/HoldingsSummaryBar.vue';
import AssetDetailDialog from './assets/AssetDetailDialog.vue';
import AssetAdminDialog from './assets/dialogs/AssetAdminDialog.vue';

import { useI18n } from '@/locales/helpers.ts';

import { useInvestmentStore } from '@/stores/investment.ts';

import { InvestmentUserAsset, type AssetInfoResponse } from '@/models/investment.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

import type { DisplayAsset } from './assets/types.ts';

const { tt } = useI18n();
const investmentStore = useInvestmentStore();

const activeTab = ref<string>('holdings');
const loading = ref<boolean>(true);
const detailDialog = ref<boolean>(false);
const selectedItem = ref<DisplayAsset | null>(null);
const showManageButton = ref<boolean>(false);
const adminDialog = ref<boolean>(false);

const watchlistAssets = ref<DisplayAsset[]>([]);
const watchlistAssetIds = ref<Set<string>>(new Set());

const holdingsCount = computed(() => investmentStore.aggregatedHoldings.length);
const watchlist = computed(() => watchlistAssets.value);

// --- Search result handling ---
function onSearchResultSelect(item: AssetInfoResponse): void {
    const existingHolding = investmentStore.aggregatedHoldings.find(a => a.assetId === item.id);
    if (existingHolding) {
        selectedItem.value = {
            assetId: existingHolding.assetId,
            assetCode: existingHolding.assetCode,
            assetName: existingHolding.assetName,
            category: existingHolding.category,
            currency: existingHolding.currency,
            market: existingHolding.market,
            industry: '',
            isHolding: true,
        };
        detailDialog.value = true;
        return;
    }

    const existingWatchlist = watchlistAssets.value.find(w => w.assetId === item.id);
    if (existingWatchlist) {
        selectedItem.value = existingWatchlist;
        detailDialog.value = true;
        return;
    }

    selectedItem.value = {
        assetId: item.id,
        assetCode: item.code,
        assetName: item.name,
        category: item.category,
        currency: item.currency,
        market: item.market,
        industry: item.industry || '',
        isHolding: false,
    };
    detailDialog.value = true;
}

function onWatchlistSelect(item: DisplayAsset): void {
    selectedItem.value = item;
    detailDialog.value = true;
}

// --- Watchlist ---
async function loadWatchlist(): Promise<void> {
    try {
        const response = await services.getUserAssets({ is_active: true, is_watchlist: true });
        const data = response.data;

        if (!data || !data.success || !data.result) {
            return;
        }

        const userAssets = InvestmentUserAsset.ofMulti(data.result);
        const holdingAssetIds = new Set(investmentStore.holdings.map(h => h.assetId));

        const items: DisplayAsset[] = [];
        const ids: Set<string> = new Set();

        for (const ua of userAssets) {
            if (!ua.asset) continue;

            const a = ua.asset;
            ids.add(a.id);

            items.push({
                assetId: a.id,
                assetCode: a.code,
                assetName: a.name,
                category: a.category,
                currency: a.currency,
                market: a.market,
                industry: a.industry || '',
                isHolding: holdingAssetIds.has(a.id),
            });
        }

        watchlistAssets.value = items;
        watchlistAssetIds.value = ids;

        // Load prices for watchlist items
        for (const item of items) {
            try {
                const result = await investmentStore.loadLatestMarketData({ assetId: item.assetId });
                if (result) {
                    item.currentPrice = result.price;
                }
            } catch {
                // Ignore individual failures
            }
        }
    } catch (error) {
        logger.error('Failed to load watchlist', error);
    }
}

function reloadWatchlist(): void {
    loadWatchlist();
}

// --- Admin check ---
async function checkAdmin(): Promise<void> {
    try {
        const response = await services.checkInvestmentAdmin();
        const data = response.data;
        if (data && data.success && data.result) {
            showManageButton.value = data.result.isAdmin;
        }
    } catch (error) {
        logger.error('Failed to check admin status', error);
    }
}

// --- Lifecycle ---
onMounted(async () => {
    loading.value = true;
    try {
        checkAdmin();
        await investmentStore.loadHoldings({ force: false });
        await loadWatchlist();
    } catch (error) {
        logger.error('Failed to load assets page data', error);
    } finally {
        loading.value = false;
    }
});
</script>

<style scoped>
.page-content {
    padding: 24px;
}

.page-body {
    display: flex;
    flex-direction: column;
    gap: 16px;
}
</style>

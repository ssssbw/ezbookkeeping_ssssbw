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

            <v-card class="flex-grow-1 d-flex flex-column">
                <div class="d-flex align-center px-4">
                    <v-tabs v-model="activeTab">
                        <v-tab value="holdings">
                            {{ tt('Holdings') }}
                            <v-chip size="small" variant="tonal" class="ml-2">{{ holdingsCount }}</v-chip>
                        </v-tab>
                        <v-tab value="watchlist">
                            {{ tt('Watchlist') }}
                            <v-chip size="small" variant="tonal" class="ml-2">{{ watchlist.length }}</v-chip>
                        </v-tab>
                    </v-tabs>
                    <v-spacer />
                    <v-btn size="small" variant="tonal" :loading="refreshing" @click="refreshPrices">
                        {{ tt('Refresh') }}
                    </v-btn>
                </div>

                <HoldingsTab v-if="activeTab === 'holdings'" :loading="loading" />

                <WatchlistTab
                    v-if="activeTab === 'watchlist'"
                    :loading="loading"
                    :watchlist="watchlist"
                    @removed="reloadWatchlist"
                />

                <template v-if="activeTab === 'holdings'">
                    <v-divider />
                    <HoldingsSummaryBar />
                </template>
            </v-card>

            <AssetAdminDialog v-model="adminDialog" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

import AssetSearchBar from './assets/AssetSearchBar.vue';
import HoldingsTab from './assets/HoldingsTab.vue';
import WatchlistTab from './assets/WatchlistTab.vue';
import HoldingsSummaryBar from './assets/HoldingsSummaryBar.vue';
import AssetAdminDialog from './assets/dialogs/AssetAdminDialog.vue';

import { useI18n } from '@/locales/helpers.ts';

import { useInvestmentStore } from '@/stores/investment.ts';

import { InvestmentUserAsset, type AssetInfoResponse } from '@/models/investment.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

import type { DisplayAsset } from './assets/types.ts';

const { tt } = useI18n();
const router = useRouter();
const investmentStore = useInvestmentStore();

const activeTab = ref<string>('holdings');
const loading = ref<boolean>(true);
const refreshing = ref<boolean>(false);
const showManageButton = ref<boolean>(false);
const adminDialog = ref<boolean>(false);

const watchlistAssets = ref<DisplayAsset[]>([]);
const watchlistAssetIds = ref<Set<string>>(new Set());

const holdingsCount = computed(() => investmentStore.aggregatedHoldings.length);
const watchlist = computed(() => watchlistAssets.value);

// --- Search result handling ---
function onSearchResultSelect(item: AssetInfoResponse): void {
    router.push(`/investment/assets/${item.id}`);
}

// --- Watchlist ---
async function loadWatchlist(): Promise<void> {
    try {
        const response = await services.getUserAssets({ is_active: true, is_watchlist: true });
        const data = response.data;

        if (!data || !data.success) {
            return;
        }

        // Handle empty list (backend returns null instead of empty array)
        if (!data.result) {
            watchlistAssets.value = [];
            watchlistAssetIds.value = new Set();
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

        // Update list first to refresh UI immediately
        watchlistAssets.value = items;
        watchlistAssetIds.value = ids;

        // Load prices in parallel (don't block UI update)
        const pricePromises = items.map(async (item) => {
            try {
                const result = await investmentStore.loadLatestMarketData({ assetId: item.assetId });
                if (result) {
                    item.currentPrice = result.price;
                    item.currentPriceDate = result.date;
                }
            } catch {
                // Ignore individual failures
            }
        });
        await Promise.allSettled(pricePromises);
    } catch (error) {
        logger.error('Failed to load watchlist', error);
    }
}

function reloadWatchlist(): void {
    loadWatchlist();
}

// --- Refresh prices ---
async function refreshPrices(): Promise<void> {
    refreshing.value = true;
    try {
        await investmentStore.refreshAllMarketData();
        await investmentStore.loadHoldings({ force: true });
        await loadWatchlist();
    } catch (error) {
        logger.error('Failed to refresh prices', error);
    } finally {
        refreshing.value = false;
    }
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
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow-y: auto;
}

.page-body {
    display: flex;
    flex-direction: column;
    gap: 16px;
    flex: 1;
    min-height: 0;
}
</style>

<style>
.page-content-container {
    height: 100%;
}
</style>

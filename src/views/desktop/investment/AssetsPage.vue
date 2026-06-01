<template>
    <div class="page-content">
        <div class="page-header">
            <h1 class="page-title">{{ tt('Asset Management') }}</h1>
        </div>
        <div class="page-body">
            <!-- Tabs and Filters Card -->
            <v-card>
                <v-card-text>
                    <v-tabs v-model="activeTab" class="mb-4">
                        <v-tab value="holdings">
                            {{ tt('Holdings') }} ({{ filteredHoldings.length }})
                        </v-tab>
                        <v-tab value="watchlist">
                            {{ tt('Watchlist') }} ({{ filteredWatchlist.length }})
                        </v-tab>
                    </v-tabs>

                    <v-row class="mb-2">
                        <v-col cols="12" md="4">
                            <v-text-field
                                v-model="searchQuery"
                                :label="tt('Search')"
                                prepend-inner-icon="mdi-magnify"
                                clearable
                                hide-details
                                density="compact"
                                variant="outlined"
                            />
                        </v-col>
                        <v-col cols="12" md="4">
                            <v-select
                                v-model="categoryFilter"
                                :label="tt('Asset Categories')"
                                :items="categoryOptions"
                                clearable
                                hide-details
                                density="compact"
                                variant="outlined"
                            />
                        </v-col>
                        <v-col cols="12" md="4">
                            <v-select
                                v-model="marketFilter"
                                :label="tt('Market')"
                                :items="marketOptions"
                                clearable
                                hide-details
                                density="compact"
                                variant="outlined"
                            />
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>

            <!-- Asset Table -->
            <v-card>
                <v-card-text>
                    <v-data-table
                        :headers="computedHeaders"
                        :items="currentTabFilteredItems"
                        :loading="loading"
                        :sort-by="sortBy"
                        :hover="true"
                        item-value="assetId"
                        @click:row="onRowClick"
                    >
                        <template #item.assetName="{ item }">
                            <span class="font-weight-medium">{{ item.assetName }}</span>
                        </template>
                        <template #item.currentPrice="{ item }">
                            <span v-if="item.currentPrice !== undefined && item.currentPrice !== null">
                                {{ formatPrice(item.currentPrice) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.marketValue="{ item }">
                            <span v-if="item.marketValue !== undefined && item.marketValue !== null" class="font-weight-medium">
                                {{ formatCurrencyValue(item.marketValue, item.currency) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.returnRate="{ item }">
                            <span v-if="item.returnRate !== undefined && item.returnRate !== null" :class="getReturnColorClass(item.returnRate)">
                                {{ formatReturnRate(item.returnRate) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.market="{ item }">
                            {{ formatMarket(item.market) }}
                        </template>
                        <template #item.category="{ item }">
                            {{ formatCategory(item.category) }}
                        </template>
                        <template #item.quantity="{ item }">
                            <span v-if="item.quantity !== undefined && item.quantity !== null">
                                {{ formatQuantity(item.quantity) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #loading>
                            <v-skeleton-loader type="table-row@10" :loading="true"></v-skeleton-loader>
                        </template>
                        <template #no-data>
                            <div class="text-center py-4 text-medium-emphasis">
                                {{ activeTab === 'holdings' ? tt('No holdings data') : tt('No watchlist data') }}
                            </div>
                        </template>
                    </v-data-table>
                </v-card-text>
            </v-card>

            <!-- Detail Dialog -->
            <v-dialog v-model="detailDialog" max-width="520">
                <v-card v-if="selectedItem">
                    <v-card-title class="d-flex align-center">
                        <span class="text-h6">{{ selectedItem.assetName }}</span>
                        <v-spacer />
                        <v-chip size="small" variant="tonal" :color="getCategoryColor(selectedItem.category)">
                            {{ formatCategory(selectedItem.category) }}
                        </v-chip>
                    </v-card-title>
                    <v-card-text>
                        <v-row dense>
                            <v-col cols="6">
                                <div class="text-caption text-medium-emphasis">{{ tt('Code') }}</div>
                                <div class="text-body-1">{{ selectedItem.assetCode }}</div>
                            </v-col>
                            <v-col cols="6">
                                <div class="text-caption text-medium-emphasis">{{ tt('Market') }}</div>
                                <div class="text-body-1">{{ formatMarket(selectedItem.market) }}</div>
                            </v-col>
                            <v-col cols="6">
                                <div class="text-caption text-medium-emphasis">{{ tt('Category') }}</div>
                                <div class="text-body-1">{{ formatCategory(selectedItem.category) }}</div>
                            </v-col>
                            <v-col cols="6">
                                <div class="text-caption text-medium-emphasis">{{ tt('Currency') }}</div>
                                <div class="text-body-1">{{ selectedItem.currency }}</div>
                            </v-col>
                        </v-row>

                        <template v-if="selectedItem.isHolding">
                            <v-divider class="my-3" />
                            <v-row dense>
                                <v-col cols="6">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Current Price') }}</div>
                                    <div class="text-body-1 font-weight-medium">
                                        {{ formatPrice(selectedItem.currentPrice) }}
                                    </div>
                                </v-col>
                                <v-col cols="6">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Holdings') }}</div>
                                    <div class="text-body-1">{{ formatQuantity(selectedItem.quantity) }}</div>
                                </v-col>
                                <v-col cols="6">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Avg Cost') }}</div>
                                    <div class="text-body-1">{{ formatPrice(selectedItem.avgCostPrice) }}</div>
                                </v-col>
                                <v-col cols="6">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Total Cost') }}</div>
                                    <div class="text-body-1">{{ formatCurrencyValue(selectedItem.totalCost, selectedItem.currency) }}</div>
                                </v-col>
                                <v-col cols="6">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Market Value') }}</div>
                                    <div class="text-body-1 font-weight-medium">{{ formatCurrencyValue(selectedItem.marketValue, selectedItem.currency) }}</div>
                                </v-col>
                                <v-col cols="6">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Unrealized P&L') }}</div>
                                    <div class="text-body-1" :class="getReturnColorClass(selectedItem.unrealizedPnl)">
                                        {{ formatCurrencyValue(selectedItem.unrealizedPnl, selectedItem.currency) }}
                                    </div>
                                </v-col>
                                <v-col cols="12">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Return Rate') }}</div>
                                    <div class="text-body-1 font-weight-bold" :class="getReturnColorClass(selectedItem.returnRate)">
                                        {{ formatReturnRate(selectedItem.returnRate) }}
                                    </div>
                                </v-col>
                            </v-row>
                        </template>

                        <template v-else>
                            <v-divider class="my-3" />
                            <v-row dense>
                                <v-col cols="12">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Current Price') }}</div>
                                    <div class="text-body-1" v-if="selectedItem.currentPrice !== undefined && selectedItem.currentPrice !== null">
                                        {{ formatPrice(selectedItem.currentPrice) }}
                                    </div>
                                    <div class="text-body-1 text-medium-emphasis" v-else>--</div>
                                </v-col>
                            </v-row>
                        </template>
                    </v-card-text>
                    <v-card-actions>
                        <v-spacer />
                        <v-btn variant="text" @click="detailDialog = false">{{ tt('Close') }}</v-btn>
                        <template v-if="!selectedItem.isHolding">
                            <v-btn color="error" variant="text" @click="removeFromWatchlist(selectedItem)">
                                {{ tt('Remove') }}
                            </v-btn>
                        </template>
                    </v-card-actions>
                </v-card>
            </v-dialog>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useInvestmentStore } from '@/stores/investment.ts';

import {
    AssetCategory,
    InvestmentMarket,
    InvestmentHolding,
    InvestmentUserAsset
} from '@/models/investment.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

const { tt } = useI18n();

const investmentStore = useInvestmentStore();

// --- Constants ---
const DIVISOR = 10000;

// --- State ---
const activeTab = ref<string>('holdings');
const searchQuery = ref<string>('');
const categoryFilter = ref<string | null>(null);
const marketFilter = ref<number | null>(null);
const loading = ref<boolean>(true);
const detailDialog = ref<boolean>(false);
const watchlistAssets = ref<DisplayAsset[]>([]);
const selectedItem = ref<DisplayAsset | null>(null);

// --- Display Asset interface ---
interface DisplayAsset {
    assetId: string;
    assetCode: string;
    assetName: string;
    category: string;
    currency: string;
    market: number;
    quantity?: number;
    avgCostPrice?: number;
    totalCost?: number;
    currentPrice?: number;
    marketValue?: number;
    unrealizedPnl?: number;
    returnRate?: number;
    isHolding: boolean;
}

// --- Category options ---
const categoryOptions = computed(() => [
    { title: tt('All'), value: null as string | null },
    { title: tt('Equity'), value: AssetCategory.Equity },
    { title: tt('Fixed Income'), value: AssetCategory.FixedIncome },
    { title: tt('Commodity'), value: AssetCategory.Commodity },
    { title: tt('Digital'), value: AssetCategory.Digital }
]);

// --- Market options ---
const marketOptions = computed(() => [
    { title: tt('All'), value: null as number | null },
    { title: tt('Market CN'), value: InvestmentMarket.CN },
    { title: tt('Market HK'), value: InvestmentMarket.HK },
    { title: tt('Market US'), value: InvestmentMarket.US }
]);

// --- Table headers ---
const holdingsHeaders = [
    { key: 'assetName', value: 'assetName', title: 'Name', sortable: true },
    { key: 'assetCode', value: 'assetCode', title: 'Code', sortable: true },
    { key: 'market', value: 'market', title: 'Market', sortable: true },
    { key: 'currentPrice', value: 'currentPrice', title: 'Current Price', sortable: true, align: 'end' as const },
    { key: 'quantity', value: 'quantity', title: 'Holdings', sortable: true, align: 'end' as const },
    { key: 'marketValue', value: 'marketValue', title: 'Total Value', sortable: true, align: 'end' as const },
    { key: 'returnRate', value: 'returnRate', title: 'Return Rate', sortable: true, align: 'end' as const }
];

const watchlistHeaders = [
    { key: 'assetName', value: 'assetName', title: 'Name', sortable: true },
    { key: 'assetCode', value: 'assetCode', title: 'Code', sortable: true },
    { key: 'market', value: 'market', title: 'Market', sortable: true },
    { key: 'category', value: 'category', title: 'Category', sortable: true },
    { key: 'currentPrice', value: 'currentPrice', title: 'Current Price', sortable: true, align: 'end' as const }
];

const computedHeaders = computed(() => {
    const headers = activeTab.value === 'holdings' ? holdingsHeaders : watchlistHeaders;
    return headers.map(h => ({
        ...h,
        title: tt(h.title)
    }));
});

const sortBy = computed(() => {
    if (activeTab.value === 'holdings') {
        return [{ key: 'marketValue', order: 'desc' as const }];
    }
    return [{ key: 'assetName', order: 'asc' as const }];
});

// --- Holdings (from store) ---
const holdingsList = computed<DisplayAsset[]>(() => {
    return investmentStore.holdings.map(h => holdingToDisplay(h));
});

function holdingToDisplay(h: InvestmentHolding): DisplayAsset {
    return {
        assetId: h.assetId,
        assetCode: h.assetCode,
        assetName: h.assetName,
        category: h.category,
        currency: h.currency,
        market: h.market,
        quantity: h.quantity,
        avgCostPrice: h.avgCostPrice,
        totalCost: h.totalCost,
        currentPrice: h.currentPrice,
        marketValue: h.marketValue,
        unrealizedPnl: h.unrealizedPnl,
        returnRate: h.returnRate,
        isHolding: true
    };
}

// --- Filtered lists ---
function filterAssets(assets: DisplayAsset[]): DisplayAsset[] {
    let result = assets;

    if (searchQuery.value.trim()) {
        const query = searchQuery.value.trim().toLowerCase();
        result = result.filter(a =>
            a.assetName.toLowerCase().includes(query) ||
            a.assetCode.toLowerCase().includes(query)
        );
    }

    if (categoryFilter.value) {
        result = result.filter(a => a.category === categoryFilter.value);
    }

    if (marketFilter.value !== null) {
        result = result.filter(a => a.market === marketFilter.value);
    }

    return result;
}

const filteredHoldings = computed<DisplayAsset[]>(() => {
    return filterAssets(holdingsList.value);
});

const filteredWatchlist = computed<DisplayAsset[]>(() => {
    return filterAssets(watchlistAssets.value);
});

const currentTabFilteredItems = computed<DisplayAsset[]>(() => {
    return activeTab.value === 'holdings' ? filteredHoldings.value : filteredWatchlist.value;
});

// --- Watchlist loading ---
async function loadWatchlist(): Promise<void> {
    try {
        const response = await services.getUserAssets({ is_active: true });
        const data = response.data;

        if (!data || !data.success || !data.result) {
            logger.error('Failed to load watchlist assets');
            return;
        }

        const userAssets = InvestmentUserAsset.ofMulti(data.result);

        const holdingAssetIds = new Set(investmentStore.holdings.map(h => h.assetId));

        const watchlistItems: DisplayAsset[] = [];

        for (const ua of userAssets) {
            if (holdingAssetIds.has(ua.assetId)) {
                continue;
            }

            if (!ua.asset) {
                continue;
            }

            const a = ua.asset;
            watchlistItems.push({
                assetId: a.id,
                assetCode: a.code,
                assetName: a.name,
                category: a.category,
                currency: a.currency,
                market: a.market,
                isHolding: false
            });
        }

        watchlistAssets.value = watchlistItems;

        loadWatchlistPrices(watchlistItems);
    } catch (error) {
        logger.error('Failed to load watchlist', error);
    }
}

async function loadWatchlistPrices(items: DisplayAsset[]): Promise<void> {
    const pricePromises = items.map(async (item) => {
        try {
            const result = await investmentStore.loadLatestMarketData({ assetId: item.assetId });
            if (result) {
                item.currentPrice = result.price;
            }
        } catch {
            // Ignore individual failures
        }
    });

    await Promise.all(pricePromises);
}

// --- Remove from watchlist ---
async function removeFromWatchlist(item: DisplayAsset): Promise<void> {
    try {
        await investmentStore.removeUserAsset({ assetId: item.assetId });
        watchlistAssets.value = watchlistAssets.value.filter(a => a.assetId !== item.assetId);
        detailDialog.value = false;
    } catch (error) {
        logger.error('Failed to remove asset from watchlist', error);
    }
}

// --- Detail dialog ---
function onRowClick(_event: Event, row: { item: DisplayAsset }): void {
    showDetail(row.item);
}

function showDetail(item: DisplayAsset): void {
    selectedItem.value = item;
    detailDialog.value = true;
}

// --- Format helpers ---
function formatPrice(value: number | undefined | null): string {
    if (value === undefined || value === null) return '--';
    const num = value / DIVISOR;
    if (num >= 1000) return num.toFixed(2);
    if (num >= 1) return num.toFixed(4);
    return num.toFixed(6);
}

function formatCurrencyValue(value: number | undefined | null, currency: string): string {
    if (value === undefined || value === null) return '--';
    const num = value / DIVISOR;
    const abs = Math.abs(num);
    let formatted: string;
    if (abs >= 1000000) {
        formatted = (abs / 1000000).toFixed(2) + 'M';
    } else if (abs >= 10000) {
        formatted = (abs / 10000).toFixed(2) + 'W';
    } else if (abs >= 1000) {
        formatted = (abs / 1000).toFixed(2) + 'K';
    } else {
        formatted = abs.toFixed(2);
    }
    const sign = num < 0 ? '-' : '';
    return sign + getCurrencySymbol(currency) + formatted;
}

function formatReturnRate(value: number | undefined | null): string {
    if (value === undefined || value === null) return '--';
    const pct = value / 100;
    const sign = pct >= 0 ? '+' : '';
    return sign + pct.toFixed(2) + '%';
}

function formatQuantity(value: number | undefined | null): string {
    if (value === undefined || value === null) return '--';
    const num = value / DIVISOR;
    if (num >= 1000000) return (num / 1000000).toFixed(2) + 'M';
    if (num >= 10000) return (num / 10000).toFixed(2) + 'W';
    if (num >= 1000) return (num / 1000).toFixed(2) + 'K';
    return num.toFixed(2);
}

function formatMarket(market: number): string {
    switch (market) {
        case InvestmentMarket.CN: return tt('Market CN');
        case InvestmentMarket.HK: return tt('Market HK');
        case InvestmentMarket.US: return tt('Market US');
        default: return market.toString();
    }
}

function formatCategory(category: string): string {
    switch (category) {
        case AssetCategory.Equity: return tt('Equity');
        case AssetCategory.FixedIncome: return tt('Fixed Income');
        case AssetCategory.Commodity: return tt('Commodity');
        case AssetCategory.Digital: return tt('Digital');
        default: return category;
    }
}

function getCurrencySymbol(currency: string): string {
    switch (currency) {
        case 'CNY': return '\u00a5';
        case 'HKD': return 'HK$';
        case 'USD': return '$';
        default: return '';
    }
}

function getReturnColorClass(value: number | undefined | null): string {
    if (value === undefined || value === null) return '';
    if (value > 0) return 'text-success';
    if (value < 0) return 'text-error';
    return '';
}

function getCategoryColor(category: string): string {
    switch (category) {
        case AssetCategory.Equity: return 'primary';
        case AssetCategory.FixedIncome: return 'success';
        case AssetCategory.Commodity: return 'warning';
        case AssetCategory.Digital: return 'purple';
        default: return 'default';
    }
}

// --- Lifecycle ---
onMounted(async () => {
    loading.value = true;
    try {
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

.page-header {
    margin-bottom: 24px;
}

.page-title {
    font-size: 24px;
    font-weight: 600;
    margin: 0;
}

.page-body {
    display: flex;
    flex-direction: column;
    gap: 24px;
}
</style>

<template>
    <div class="page-content">
        <div class="page-body">
            <!-- Search with dropdown -->
            <div class="position-relative">
                <v-text-field
                    v-model="searchKeyword"
                    :label="tt('Search')"
                    :placeholder="tt('Search assets')"
                    prepend-inner-icon="mdi-magnify"
                    clearable
                    hide-details
                    density="default"
                    variant="outlined"
                    :loading="searchLoading"
                    @update:model-value="onSearchInput"
                    @click:clear="clearSearch"
                />
                <!-- Search Results Dropdown -->
                <v-menu
                    v-model="showSearchResults"
                    :close-on-content-click="false"
                    max-height="400"
                    min-width="400"
                    location="bottom"
                >
                    <template #activator="{ props }">
                        <div v-bind="props" class="search-anchor"></div>
                    </template>
                    <v-list density="compact" lines="two">
                        <v-list-item
                            v-for="item in searchResults"
                            :key="item.id"
                        >
                            <template #title>
                                <span class="text-body-2 font-weight-medium">{{ item.code }}</span>
                                <span class="text-body-2 ms-2">{{ item.name }}</span>
                            </template>
                            <template #subtitle>
                                <span class="text-caption text-medium-emphasis">
                                    {{ formatCategory(item.category) }} · {{ formatMarket(item.market) }}
                                </span>
                            </template>
                            <template #prepend>
                                <v-chip
                                    v-if="isInHoldings(item.id)"
                                    size="x-small"
                                    color="success"
                                    variant="tonal"
                                    class="me-2"
                                >
                                    {{ tt('Holdings') }}
                                </v-chip>
                                <v-chip
                                    v-else-if="isInWatchlist(item.id)"
                                    size="x-small"
                                    color="primary"
                                    variant="tonal"
                                    class="me-2"
                                >
                                    {{ tt('Watchlist') }}
                                </v-chip>
                            </template>
                            <template #append>
                                <v-btn
                                    v-if="!isInHoldings(item.id) && !isInWatchlist(item.id)"
                                    size="x-small"
                                    variant="tonal"
                                    color="primary"
                                    @click.stop="addAssetToWatchlist(item)"
                                >
                                    {{ tt('Add to Watchlist') }}
                                </v-btn>
                                <v-btn
                                    size="x-small"
                                    variant="tonal"
                                    color="success"
                                    class="ms-1"
                                    @click.stop="openTransactionDialog(item); showSearchResults = false"
                                >
                                    {{ tt('Buy') }}
                                </v-btn>
                            </template>
                        </v-list-item>
                    </v-list>
                </v-menu>
            </div>

            <!-- Filters -->
            <div class="d-flex flex-wrap ga-2">
                <v-chip-group v-model="categoryFilter" mandatory>
                    <v-chip value="" filter variant="outlined" size="small">{{ tt('All') }}</v-chip>
                    <v-chip :value="AssetCategory.Equity" filter variant="outlined" size="small">{{ tt('Equity') }}</v-chip>
                    <v-chip :value="AssetCategory.FixedIncome" filter variant="outlined" size="small">{{ tt('Fixed Income') }}</v-chip>
                    <v-chip :value="AssetCategory.Commodity" filter variant="outlined" size="small">{{ tt('Commodity') }}</v-chip>
                    <v-chip :value="AssetCategory.Digital" filter variant="outlined" size="small">{{ tt('Digital') }}</v-chip>
                </v-chip-group>
                <v-divider vertical class="mx-1" />
                <v-chip-group v-model="marketFilter" mandatory>
                    <v-chip :value="null" filter variant="outlined" size="small">{{ tt('All Markets') }}</v-chip>
                    <v-chip :value="InvestmentMarket.CN" filter variant="outlined" size="small">{{ tt('Market CN') }}</v-chip>
                    <v-chip :value="InvestmentMarket.HK" filter variant="outlined" size="small">{{ tt('Market HK') }}</v-chip>
                    <v-chip :value="InvestmentMarket.US" filter variant="outlined" size="small">{{ tt('Market US') }}</v-chip>
                </v-chip-group>
            </div>

            <!-- Main Card: Tabs + Table -->
            <v-card>
                <v-card-title class="d-flex align-center py-2">
                    <v-tabs v-model="activeTab" density="compact" class="flex-grow-1">
                        <v-tab value="watchlist">
                            {{ tt('Watchlist') }} ({{ watchlistCount }})
                        </v-tab>
                        <v-tab value="holdings">
                            {{ tt('Holdings') }} ({{ holdingsCount }})
                        </v-tab>
                    </v-tabs>
                    <v-btn
                        v-if="isAdmin"
                        variant="tonal"
                        color="primary"
                        size="small"
                        class="ms-2"
                        @click="adminDialog = true"
                    >
                        {{ tt('Asset Management') }}
                    </v-btn>
                </v-card-title>

                <v-divider />

                <v-card-text class="pt-2">
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
                        <template #bottom>
                            <div class="title-and-toolbar d-flex align-center text-no-wrap mt-2" v-if="currentTabFilteredItems">
                                <span class="text-body-2 text-medium-emphasis">
                                    {{ tt('format.misc.selectedCount', { count: formatNumberToLocalizedNumerals(currentTabFilteredItems.length), totalCount: formatNumberToLocalizedNumerals(currentTabFilteredItems.length) }) }}
                                </span>
                                <v-spacer v-if="currentTabFilteredItems.length > 10" />
                                <span v-if="currentTabFilteredItems.length > 10" class="text-body-2 text-medium-emphasis">{{ tt('Transactions Per Page') }}</span>
                                <v-select class="ms-2" density="compact" max-width="100"
                                          item-title="name" item-value="value"
                                          :items="getTablePageOptions(currentTabFilteredItems.length)"
                                          v-model="countPerPage"
                                          v-if="currentTabFilteredItems.length > 10"
                                />
                                <pagination-buttons density="compact"
                                                    :totalPageCount="totalPageCount"
                                                    v-model="currentPage"
                                                    v-if="currentTabFilteredItems.length > 10"></pagination-buttons>
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
                                <div class="text-caption text-medium-emphasis">{{ tt('Asset Categories') }}</div>
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
                                    <div class="text-caption text-medium-emphasis">{{ tt('Quantity') }}</div>
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

            <!-- Buy/Sell Dialog -->
            <v-dialog v-model="transactionDialog" max-width="520" persistent>
                <v-card v-if="selectedTransactionAsset">
                    <v-card-title class="d-flex align-center">
                        <span class="text-h6">{{ selectedTransactionAsset.name || selectedTransactionAsset.code }}</span>
                        <v-spacer />
                        <span class="text-body-2 text-medium-emphasis">{{ selectedTransactionAsset.code }}</span>
                    </v-card-title>
                    <v-card-text>
                        <v-form @submit.prevent="submitTransaction">
                            <v-row dense>
                                <v-col cols="12">
                                    <v-select
                                        v-model="transactionForm.accountId"
                                        :label="tt('Account')"
                                        :items="investmentAccountItems"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        required
                                    />
                                </v-col>
                                <v-col cols="12">
                                    <v-select
                                        v-model="transactionForm.type"
                                        :label="tt('Transaction Type')"
                                        :items="transactionTypeOptions"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        required
                                    />
                                </v-col>
                                <v-col cols="12" md="6">
                                    <v-text-field
                                        v-model.number="transactionForm.amount"
                                        :label="tt('Amount')"
                                        type="number"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        step="0.01"
                                        min="0"
                                    />
                                </v-col>
                                <v-col cols="12" md="6">
                                    <v-text-field
                                        v-model.number="transactionForm.quantity"
                                        :label="tt('Quantity')"
                                        type="number"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        step="0.01"
                                        min="0"
                                    />
                                </v-col>
                                <v-col cols="12" md="6">
                                    <v-text-field
                                        v-model.number="transactionForm.price"
                                        :label="tt('Current Price')"
                                        type="number"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        step="0.0001"
                                        min="0"
                                    />
                                </v-col>
                                <v-col cols="12" md="6">
                                    <v-text-field
                                        v-model.number="transactionForm.fee"
                                        :label="tt('Service Charge')"
                                        type="number"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        step="0.01"
                                        min="0"
                                    />
                                </v-col>
                                <v-col cols="12">
                                    <v-text-field
                                        v-model="transactionForm.tradeTime"
                                        :label="tt('Transaction Time')"
                                        type="datetime-local"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                    />
                                </v-col>
                                <v-col cols="12">
                                    <v-textarea
                                        v-model="transactionForm.comment"
                                        :label="tt('Comment')"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        rows="2"
                                        :placeholder="tt('Comment')"
                                    />
                                </v-col>
                            </v-row>
                        </v-form>
                    </v-card-text>
                    <v-card-actions>
                        <v-spacer />
                        <v-btn variant="text" :disabled="transactionSubmitting" @click="transactionDialog = false">
                            {{ tt('Cancel') }}
                        </v-btn>
                        <v-btn
                            color="primary"
                            variant="elevated"
                            :disabled="!isTransactionFormValid || transactionSubmitting"
                            :loading="transactionSubmitting"
                            @click="submitTransaction"
                        >
                            {{ tt('Save') }}
                        </v-btn>
                    </v-card-actions>
                </v-card>
            </v-dialog>

            <!-- Admin Dialog -->
            <v-dialog v-model="adminDialog" max-width="600">
                <v-card>
                    <v-card-title class="d-flex align-center">
                        <span class="text-h6">{{ tt('Asset Management') }}</span>
                        <v-spacer />
                        <v-btn variant="text" icon="mdi-close" density="compact" @click="adminDialog = false" />
                    </v-card-title>
                    <v-card-text>
                        <div class="text-body-2 text-medium-emphasis">
                            {{ tt('Asset Management') }}
                        </div>
                    </v-card-text>
                    <v-card-actions>
                        <v-spacer />
                        <v-btn variant="text" @click="adminDialog = false">{{ tt('Close') }}</v-btn>
                    </v-card-actions>
                </v-card>
            </v-dialog>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import PaginationButtons from '@/components/desktop/PaginationButtons.vue';

import { useInvestmentStore } from '@/stores/investment.ts';
import { useAccountsStore } from '@/stores/account.ts';

import { AccountCategory } from '@/core/account.ts';

import {
    AssetCategory,
    InvestmentMarket,
    InvestmentTransactionType,
    InvestmentHolding,
    InvestmentUserAsset,
    type AssetInfoResponse
} from '@/models/investment.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';
import { getCurrentUnixTime, getTimezoneOffsetMinutes } from '@/lib/datetime.ts';

const { tt, formatNumberToLocalizedNumerals } = useI18n();

const investmentStore = useInvestmentStore();
const accountsStore = useAccountsStore();

// --- Constants ---
const DIVISOR = 10000;

// --- State ---
const activeTab = ref<string>('holdings');
const loading = ref<boolean>(true);
const detailDialog = ref<boolean>(false);
const selectedItem = ref<DisplayAsset | null>(null);
const isAdmin = ref<boolean>(false);
const adminDialog = ref<boolean>(false);

// --- Filter State ---
const categoryFilter = ref<string>('');
const marketFilter = ref<number | null>(null);

// --- Search State ---
const searchKeyword = ref<string>('');
const searchResults = ref<AssetInfoResponse[]>([]);
const searchLoading = ref<boolean>(false);
const showSearchResults = ref<boolean>(false);
let searchTimer: ReturnType<typeof setTimeout> | null = null;

// --- Watchlist ---
const watchlistAssets = ref<DisplayAsset[]>([]);
const watchlistAssetIds = ref<Set<string>>(new Set());

// --- Transaction Dialog ---
const transactionDialog = ref<boolean>(false);
const transactionSubmitting = ref<boolean>(false);
const selectedTransactionAsset = ref<AssetInfoResponse | null>(null);
const transactionForm = ref({
    accountId: '',
    type: InvestmentTransactionType.Buy,
    amount: 0,
    quantity: 0,
    price: 0,
    fee: 0,
    tradeTime: '',
    comment: ''
});

// --- Pagination ---
const countPerPage = ref<number>(10);
const currentPage = ref<number>(1);

const totalPageCount = computed<number>(() => {
    return Math.ceil(currentTabFilteredItems.value.length / countPerPage.value);
});

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

// --- Pagination options ---
function getTablePageOptions(linesCount: number): { value: number; name: string }[] {
    const pageOptions: { value: number; name: string }[] = [];

    if (!linesCount || linesCount < 1) {
        pageOptions.push({ value: -1, name: tt('All') });
        return pageOptions;
    }

    const availableCountPerPage = [5, 10, 15, 20, 25, 30, 50];

    for (const count of availableCountPerPage) {
        if (linesCount < count) break;
        pageOptions.push({ value: count, name: formatNumberToLocalizedNumerals(count) });
    }

    pageOptions.push({ value: -1, name: tt('All') });
    return pageOptions;
}

// --- Investment account options for transaction form ---
const investmentAccountItems = computed<{ title: string; value: string }[]>(() => {
    const items: { title: string; value: string }[] = [];
    for (const account of accountsStore.allPlainAccounts) {
        if (account.category === AccountCategory.InvestmentAccount.type) {
            items.push({
                title: account.name,
                value: account.id
            });
        }
    }
    return items;
});

// --- Transaction type options ---
const transactionTypeOptions = computed(() => [
    { title: tt('Total Buys'), value: InvestmentTransactionType.Buy },
    { title: tt('Total Sells'), value: InvestmentTransactionType.Sell }
]);

const isTransactionFormValid = computed<boolean>(() => {
    return !!transactionForm.value.accountId && transactionForm.value.amount > 0;
});

// --- Table headers ---
const holdingsHeaders = [
    { key: 'assetName', value: 'assetName', title: 'Name', sortable: true },
    { key: 'assetCode', value: 'assetCode', title: 'Code', sortable: true },
    { key: 'market', value: 'market', title: 'Market', sortable: true },
    { key: 'currentPrice', value: 'currentPrice', title: 'Current Price', sortable: true, align: 'end' as const },
    { key: 'quantity', value: 'quantity', title: 'Quantity', sortable: true, align: 'end' as const },
    { key: 'marketValue', value: 'marketValue', title: 'Total Value', sortable: true, align: 'end' as const },
    { key: 'returnRate', value: 'returnRate', title: 'Return Rate', sortable: true, align: 'end' as const }
];

const watchlistHeaders = [
    { key: 'assetName', value: 'assetName', title: 'Name', sortable: true },
    { key: 'assetCode', value: 'assetCode', title: 'Code', sortable: true },
    { key: 'market', value: 'market', title: 'Market', sortable: true },
    { key: 'category', value: 'category', title: 'Asset Categories', sortable: true },
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

const holdingsCount = computed<number>(() => {
    return investmentStore.holdings.length;
});

const watchlistCount = computed<number>(() => {
    return watchlistAssets.value.length;
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
function filterByCategoryAndMarket(items: DisplayAsset[]): DisplayAsset[] {
    let result = items;

    if (categoryFilter.value) {
        result = result.filter(a => a.category === categoryFilter.value);
    }

    if (marketFilter.value !== null) {
        result = result.filter(a => a.market === marketFilter.value);
    }

    return result;
}

const filteredHoldings = computed<DisplayAsset[]>(() => {
    return filterByCategoryAndMarket(holdingsList.value);
});

const filteredWatchlist = computed<DisplayAsset[]>(() => {
    return filterByCategoryAndMarket(watchlistAssets.value);
});

const currentTabFilteredItems = computed<DisplayAsset[]>(() => {
    return activeTab.value === 'holdings' ? filteredHoldings.value : filteredWatchlist.value;
});

// --- Search ---
function clearSearch(): void {
    searchKeyword.value = '';
    searchResults.value = [];
    showSearchResults.value = false;
}

function onSearchInput(): void {
    if (searchTimer) {
        clearTimeout(searchTimer);
    }

    const keyword = searchKeyword.value.trim();
    if (!keyword) {
        searchResults.value = [];
        showSearchResults.value = false;
        return;
    }

    searchTimer = setTimeout(() => {
        performSearch(keyword);
    }, 300);
}

async function performSearch(keyword: string): Promise<void> {
    searchLoading.value = true;
    try {
        const response = await services.searchAssets({ keyword, limit: 10 });
        const data = response.data;

        if (!data || !data.success || !data.result) {
            searchResults.value = [];
            showSearchResults.value = false;
            return;
        }

        searchResults.value = data.result;
        showSearchResults.value = data.result.length > 0;
    } catch (error) {
        logger.error('Failed to search assets', error);
        searchResults.value = [];
        showSearchResults.value = false;
    } finally {
        searchLoading.value = false;
    }
}

function isInHoldings(assetId: string): boolean {
    return !!investmentStore.holdingsMap[assetId];
}

function isInWatchlist(assetId: string): boolean {
    return watchlistAssetIds.value.has(assetId);
}

// --- Watchlist loading ---
async function loadWatchlist(): Promise<void> {
    try {
        const response = await services.getUserAssets({ is_watchlist: true });
        const data = response.data;

        if (!data || !data.success || !data.result) {
            logger.error('Failed to load watchlist assets');
            return;
        }

        const userAssets = InvestmentUserAsset.ofMulti(data.result);

        const holdingAssetIds = new Set(investmentStore.holdings.map(h => h.assetId));

        const watchlistItems: DisplayAsset[] = [];
        const ids: Set<string> = new Set();

        for (const ua of userAssets) {
            if (!ua.asset) {
                continue;
            }

            const a = ua.asset;

            ids.add(a.id);

            watchlistItems.push({
                assetId: a.id,
                assetCode: a.code,
                assetName: a.name,
                category: a.category,
                currency: a.currency,
                market: a.market,
                isHolding: holdingAssetIds.has(a.id)
            });
        }

        watchlistAssets.value = watchlistItems;
        watchlistAssetIds.value = ids;

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

async function addAssetToWatchlist(asset: AssetInfoResponse): Promise<void> {
    try {
        await investmentStore.addUserAsset({ assetId: asset.id });

        watchlistAssetIds.value.add(asset.id);

        searchResults.value = searchResults.value.map(r =>
            r.id === asset.id ? { ...r } : r
        );

        await loadWatchlist();
    } catch (error) {
        logger.error('Failed to add asset to watchlist', error);
    }
}

// --- Remove from watchlist ---
async function removeFromWatchlist(item: DisplayAsset): Promise<void> {
    try {
        await investmentStore.removeUserAsset({ assetId: item.assetId });
        watchlistAssets.value = watchlistAssets.value.filter(a => a.assetId !== item.assetId);
        watchlistAssetIds.value.delete(item.assetId);
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

// --- Transaction Dialog ---
function openTransactionDialog(asset: AssetInfoResponse): void {
    selectedTransactionAsset.value = asset;

    const now = getCurrentUnixTime();
    const dt = new Date(now * 1000);
    const offsetMinutes = dt.getTimezoneOffset();
    const localISO = new Date(now * 1000 - offsetMinutes * 60000).toISOString().slice(0, 16);

    transactionForm.value = {
        accountId: investmentAccountItems.value.length > 0 ? investmentAccountItems.value[0]!.value : '',
        type: InvestmentTransactionType.Buy,
        amount: 0,
        quantity: 0,
        price: 0,
        fee: 0,
        tradeTime: localISO,
        comment: ''
    };

    transactionDialog.value = true;
}

async function submitTransaction(): Promise<void> {
    if (!selectedTransactionAsset.value || !isTransactionFormValid.value) {
        return;
    }

    transactionSubmitting.value = true;
    try {
        const form = transactionForm.value;
        const tradeTime = new Date(form.tradeTime).getTime() / 1000;

        const request = {
            assetId: selectedTransactionAsset.value.id,
            accountId: form.accountId,
            type: form.type,
            tradeTime: Math.floor(tradeTime),
            quantity: form.quantity > 0 ? Math.round(form.quantity * DIVISOR) : undefined,
            price: form.price > 0 ? Math.round(form.price * DIVISOR) : undefined,
            amount: Math.round(form.amount * DIVISOR),
            fee: form.fee > 0 ? Math.round(form.fee * DIVISOR) : undefined,
            utcOffset: getTimezoneOffsetMinutes(Math.floor(tradeTime)),
            comment: form.comment || undefined
        };

        await services.addInvestmentTransaction(request);

        transactionDialog.value = false;
        selectedTransactionAsset.value = null;
    } catch (error) {
        logger.error('Failed to create investment transaction', error);
    } finally {
        transactionSubmitting.value = false;
    }
}

// --- Admin Check ---
async function checkAdmin(): Promise<void> {
    try {
        const response = await services.checkInvestmentAdmin();
        const data = response.data;

        if (data && data.success && data.result) {
            isAdmin.value = (data.result as { isAdmin: boolean }).isAdmin || false;
        }
    } catch (error) {
        // Not an admin or API unavailable, leave isAdmin as false
        logger.debug('Failed to check investment admin status', error);
    }
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

// --- Watch for tab change to reset page ---
watch(activeTab, () => {
    currentPage.value = 1;
});

// --- Lifecycle ---
onMounted(async () => {
    loading.value = true;
    try {
        checkAdmin();
        await accountsStore.loadAllAccounts({ force: false });
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
    gap: 16px;
}

.position-relative {
    position: relative;
}

.search-anchor {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
}

:deep(.v-overlay__content) {
    border: 1px solid rgba(0, 0, 0, 0.12);
    border-radius: 4px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
</style>

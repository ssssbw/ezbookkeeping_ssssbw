<template>
    <div class="page-content">
        <div class="page-body">
            <!-- Top Search Bar -->
            <div class="search-section">
                <div class="search-wrapper">
                    <v-text-field
                        v-model="searchQuery"
                        :placeholder="tt('Search assets by name or code')"
                        prepend-inner-icon="mdi-magnify"
                        clearable
                        hide-details
                        density="comfortable"
                        variant="outlined"
                        class="search-input"
                        @focus="onSearchFocus"
                        @blur="onSearchBlur"
                        @update:model-value="onSearchInput"
                    />
                    <!-- Search Results Overlay -->
                    <v-card
                        v-if="showSearchResults && searchResults.length > 0"
                        class="search-results-overlay"
                        elevation="8"
                    >
                        <v-list density="compact" lines="two">
                            <v-list-item
                                v-for="item in searchResults"
                                :key="item.id"
                                @mousedown.prevent="onSearchResultClick(item)"
                            >
                                <template #prepend>
                                    <v-avatar size="32" :color="getCategoryColor(item.category)" variant="tonal">
                                        <span class="text-caption font-weight-bold">{{ item.code.substring(0, 2) }}</span>
                                    </v-avatar>
                                </template>
                                <v-list-item-title class="d-flex align-center">
                                    <span class="text-body-2 font-weight-medium">{{ item.code }}</span>
                                    <span class="text-body-2 text-medium-emphasis ml-2">{{ item.name }}</span>
                                </v-list-item-title>
                                <v-list-item-subtitle>
                                    <span class="text-caption text-medium-emphasis">{{ formatMarket(item.market) }}</span>
                                    <span v-if="item.currentPrice !== undefined && item.currentPrice !== null" class="text-caption ml-2">
                                        {{ formatPrice(item.currentPrice) }}
                                    </span>
                                    <span v-if="item.changeRate !== undefined && item.changeRate !== null" class="text-caption ml-2" :class="getReturnColorClass(item.changeRate)">
                                        {{ formatReturnRate(item.changeRate) }}
                                    </span>
                                </v-list-item-subtitle>
                                <template #append>
                                    <div class="d-flex ga-1">
                                        <v-btn
                                            v-if="!isInWatchlist(item.id)"
                                            size="x-small"
                                            variant="tonal"
                                            color="primary"
                                            @mousedown.prevent.stop="addToWatchlist(item)"
                                        >
                                            +{{ tt('Watchlist') }}
                                        </v-btn>
                                        <v-btn
                                            size="x-small"
                                            variant="tonal"
                                            color="success"
                                            @mousedown.prevent.stop="onBuyClick(item)"
                                        >
                                            {{ tt('Buy') }}
                                        </v-btn>
                                    </div>
                                </template>
                            </v-list-item>
                        </v-list>
                        <v-divider />
                        <div class="text-center py-2">
                            <v-btn
                                variant="text"
                                size="small"
                                color="primary"
                                @mousedown.prevent="onViewAllResults"
                            >
                                {{ tt('View all results') }} &gt;
                            </v-btn>
                        </div>
                    </v-card>
                </div>
                <v-btn
                    v-if="showManageButton"
                    variant="outlined"
                    color="primary"
                    class="manage-btn ml-3"
                    @click="onManageClick"
                >
                    {{ tt('Asset Management') }}
                </v-btn>
            </div>

            <!-- Tabs -->
            <v-card>
                <v-tabs v-model="activeTab" class="px-4">
                    <v-tab value="holdings">
                        {{ tt('Holdings') }}
                        <span class="text-medium-emphasis ml-1">({{ filteredHoldings.length }})</span>
                    </v-tab>
                    <v-tab value="watchlist">
                        {{ tt('Watchlist') }}
                        <span class="text-medium-emphasis ml-1">({{ filteredWatchlist.length }})</span>
                    </v-tab>
                </v-tabs>
            </v-card>

            <!-- Holdings Table -->
            <v-card v-if="activeTab === 'holdings'">
                <v-card-text class="pa-0">
                    <v-data-table
                        :headers="holdingsHeaders"
                        :items="filteredHoldings"
                        :loading="loading"
                        :sort-by="holdingsSortBy"
                        :hover="true"
                        item-value="assetId"
                        class="holdings-table"
                        @click:row="onRowClick"
                    >
                        <template #item.assetCode="{ item }">
                            <span class="text-body-2 font-weight-medium">{{ item.assetCode }}</span>
                        </template>
                        <template #item.assetName="{ item }">
                            <span class="text-body-2">{{ item.assetName }}</span>
                        </template>
                        <template #item.market="{ item }">
                            <span class="text-body-2">{{ formatMarket(item.market) }}</span>
                        </template>
                        <template #item.quantity="{ item }">
                            <span v-if="item.quantity !== undefined && item.quantity !== null" class="text-body-2">
                                {{ formatQuantity(item.quantity) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.currentPrice="{ item }">
                            <span v-if="item.currentPrice !== undefined && item.currentPrice !== null" class="text-body-2">
                                {{ formatPrice(item.currentPrice) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.marketValue="{ item }">
                            <span v-if="item.marketValue !== undefined && item.marketValue !== null" class="text-body-2 font-weight-medium">
                                {{ formatCurrencyValue(item.marketValue, item.currency) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.totalCost="{ item }">
                            <span v-if="item.totalCost !== undefined && item.totalCost !== null" class="text-body-2">
                                {{ formatCurrencyValue(item.totalCost, item.currency) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.unrealizedPnl="{ item }">
                            <span v-if="item.unrealizedPnl !== undefined && item.unrealizedPnl !== null" class="text-body-2" :class="getReturnColorClass(item.unrealizedPnl)">
                                {{ formatCurrencyValue(item.unrealizedPnl, item.currency) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.returnRate="{ item }">
                            <span v-if="item.returnRate !== undefined && item.returnRate !== null" class="text-body-2" :class="getReturnColorClass(item.returnRate)">
                                {{ formatReturnRate(item.returnRate) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.actions="{ item }">
                            <div class="d-flex ga-1">
                                <v-btn size="x-small" variant="tonal" color="success" @click.stop="onBuyClick(item)">
                                    {{ tt('Buy') }}
                                </v-btn>
                                <v-btn size="x-small" variant="tonal" color="error" @click.stop="onSellClick(item)">
                                    {{ tt('Sell') }}
                                </v-btn>
                                <v-btn size="x-small" variant="text" icon="mdi-dots-horizontal" @click.stop="showDetail(item)" />
                            </div>
                        </template>
                        <template #loading>
                            <v-skeleton-loader type="table-row@10" :loading="true" />
                        </template>
                        <template #no-data>
                            <div class="text-center py-6 text-medium-emphasis">
                                {{ tt('No holdings data') }}
                            </div>
                        </template>
                    </v-data-table>
                </v-card-text>

                <!-- Holdings Summary Footer -->
                <v-divider />
                <v-card-text class="py-3 px-6">
                    <div class="d-flex flex-wrap align-center ga-6">
                        <div class="summary-item">
                            <span class="text-caption text-medium-emphasis">{{ tt('Market Value') }}(USD)</span>
                            <span class="text-body-1 font-weight-bold ml-2">{{ formatCurrencyValue(holdingsSummary.totalMarketValue, 'USD') }}</span>
                        </div>
                        <div class="summary-item">
                            <span class="text-caption text-medium-emphasis">{{ tt('Unrealized P&L') }}(USD)</span>
                            <span class="text-body-1 font-weight-bold ml-2" :class="getReturnColorClass(holdingsSummary.totalUnrealizedPnl)">
                                {{ formatCurrencyValue(holdingsSummary.totalUnrealizedPnl, 'USD') }}
                            </span>
                        </div>
                        <div class="summary-item">
                            <span class="text-caption text-medium-emphasis">{{ tt('Return Rate') }}</span>
                            <span class="text-body-1 font-weight-bold ml-2" :class="getReturnColorClass(holdingsSummary.totalReturnRate)">
                                {{ formatReturnRate(holdingsSummary.totalReturnRate) }}
                            </span>
                        </div>
                        <div class="summary-item">
                            <span class="text-caption text-medium-emphasis">{{ tt('Today\'s P&L') }}</span>
                            <span class="text-body-1 font-weight-bold ml-2" :class="getReturnColorClass(holdingsSummary.todayPnl)">
                                {{ formatCurrencyValue(holdingsSummary.todayPnl, 'USD') }}
                            </span>
                        </div>
                        <div class="summary-item">
                            <span class="text-caption text-medium-emphasis">{{ tt('Today\'s Return Rate') }}</span>
                            <span class="text-body-1 font-weight-bold ml-2" :class="getReturnColorClass(holdingsSummary.todayReturnRate)">
                                {{ formatReturnRate(holdingsSummary.todayReturnRate) }}
                            </span>
                        </div>
                    </div>
                </v-card-text>
            </v-card>

            <!-- Watchlist Table -->
            <v-card v-if="activeTab === 'watchlist'">
                <v-card-text class="pa-0">
                    <v-data-table
                        :headers="watchlistHeaders"
                        :items="filteredWatchlist"
                        :loading="loading"
                        :sort-by="watchlistSortBy"
                        :hover="true"
                        item-value="assetId"
                        class="watchlist-table"
                        @click:row="onRowClick"
                    >
                        <template #item.assetCode="{ item }">
                            <span class="text-body-2 font-weight-medium">{{ item.assetCode }}</span>
                        </template>
                        <template #item.assetName="{ item }">
                            <span class="text-body-2">{{ item.assetName }}</span>
                        </template>
                        <template #item.market="{ item }">
                            <span class="text-body-2">{{ formatMarket(item.market) }}</span>
                        </template>
                        <template #item.currentPrice="{ item }">
                            <span v-if="item.currentPrice !== undefined && item.currentPrice !== null" class="text-body-2">
                                {{ formatPrice(item.currentPrice) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.changeRate="{ item }">
                            <span v-if="item.changeRate !== undefined && item.changeRate !== null" class="text-body-2" :class="getReturnColorClass(item.changeRate)">
                                {{ formatReturnRate(item.changeRate) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.actions="{ item }">
                            <div class="d-flex ga-1">
                                <v-btn size="x-small" variant="tonal" color="success" @click.stop="onBuyClick(item)">
                                    {{ tt('Buy') }}
                                </v-btn>
                                <v-btn size="x-small" variant="tonal" color="error" @click.stop="removeFromWatchlist(item)">
                                    {{ tt('Remove') }}
                                </v-btn>
                                <v-btn size="x-small" variant="text" icon="mdi-dots-horizontal" @click.stop="showDetail(item)" />
                            </div>
                        </template>
                        <template #loading>
                            <v-skeleton-loader type="table-row@10" :loading="true" />
                        </template>
                        <template #no-data>
                            <div class="text-center py-6 text-medium-emphasis">
                                {{ tt('No watchlist data') }}
                            </div>
                        </template>
                    </v-data-table>
                </v-card-text>
            </v-card>

            <!-- Data Delay Notice -->
            <div class="text-center text-caption text-medium-emphasis py-2">
                {{ tt('Data delayed 15 minutes') }}
            </div>

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

const { tt } = useI18n();

const investmentStore = useInvestmentStore();
const accountsStore = useAccountsStore();

// --- Constants ---
const DIVISOR = 10000;

// --- State ---
const activeTab = ref<string>('holdings');
const searchQuery = ref<string>('');
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
    changeRate?: number;
    isHolding: boolean;
}

interface SearchResultItem {
    id: string;
    code: string;
    name: string;
    category: string;
    currency: string;
    market: number;
    currentPrice?: number;
    changeRate?: number;
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

// --- Table headers ---
const holdingsHeaders = computed(() => [
    { key: 'assetCode', title: tt('Code'), sortable: true },
    { key: 'assetName', title: tt('Asset Name'), sortable: true },
    { key: 'market', title: tt('Market'), sortable: true },
    { key: 'quantity', title: tt('Quantity'), sortable: true, align: 'end' as const },
    { key: 'currentPrice', title: tt('Current Price'), sortable: true, align: 'end' as const },
    { key: 'marketValue', title: tt('Market Value'), sortable: true, align: 'end' as const },
    { key: 'totalCost', title: tt('Total Cost'), sortable: true, align: 'end' as const },
    { key: 'unrealizedPnl', title: tt('Unrealized P&L'), sortable: true, align: 'end' as const },
    { key: 'returnRate', title: tt('Return Rate'), sortable: true, align: 'end' as const },
    { key: 'actions', title: tt('Action'), sortable: false, align: 'center' as const, width: '160' }
]);

const watchlistHeaders = computed(() => [
    { key: 'assetCode', title: tt('Code'), sortable: true },
    { key: 'assetName', title: tt('Asset Name'), sortable: true },
    { key: 'market', title: tt('Market'), sortable: true },
    { key: 'currentPrice', title: tt('Current Price'), sortable: true, align: 'end' as const },
    { key: 'changeRate', title: tt('Change Rate'), sortable: true, align: 'end' as const },
    { key: 'actions', title: tt('Action'), sortable: false, align: 'center' as const, width: '160' }
]);

const holdingsSortBy = computed(() => [
    { key: 'marketValue', order: 'desc' as const }
]);

const watchlistSortBy = computed(() => [
    { key: 'assetName', order: 'asc' as const }
]);

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

// --- Holdings Summary ---
const holdingsSummary = computed(() => {
    const holdings = filteredHoldings.value;
    let totalMarketValue = 0;
    let totalUnrealizedPnl = 0;
    let totalCost = 0;

    for (const h of holdings) {
        if (h.marketValue !== undefined && h.marketValue !== null) {
            totalMarketValue += h.marketValue;
        }
        if (h.unrealizedPnl !== undefined && h.unrealizedPnl !== null) {
            totalUnrealizedPnl += h.unrealizedPnl;
        }
        if (h.totalCost !== undefined && h.totalCost !== null) {
            totalCost += h.totalCost;
        }
    }

    const totalReturnRate = totalCost !== 0
        ? Math.round((totalUnrealizedPnl / totalCost) * 10000)
        : 0;

    return {
        totalMarketValue,
        totalUnrealizedPnl,
        totalReturnRate,
        todayPnl: 0,
        todayReturnRate: 0
    };
});

// --- Filtered lists ---
function filterByCategoryAndMarket(items: DisplayAsset[]): DisplayAsset[] {
    let result = items;

    return result;
}

const filteredHoldings = computed<DisplayAsset[]>(() => {
    return filterByCategoryAndMarket(holdingsList.value);
});

const filteredWatchlist = computed<DisplayAsset[]>(() => {
    return filterByCategoryAndMarket(watchlistAssets.value);
});

// --- Search ---
let searchTimer: ReturnType<typeof setTimeout> | null = null;

function onSearchFocus(): void {
    if (searchQuery.value.trim() && searchResults.value.length > 0) {
        showSearchResults.value = true;
    }
}

function onSearchBlur(): void {
    // Delay hiding to allow click on results
    setTimeout(() => {
        showSearchResults.value = false;
    }, 200);
}

function onSearchInput(value: string): void {
    if (!value || !value.trim()) {
        searchResults.value = [];
        showSearchResults.value = false;
        return;
    }

    if (searchTimer) {
        clearTimeout(searchTimer);
    }

    searchTimer = setTimeout(() => {
        performSearch(value.trim());
    }, 300);
}

async function performSearch(keyword: string): Promise<void> {
    searchLoading.value = true;
    try {
        const response = await services.searchAssets({ keyword, limit: 8 });
        const data = response.data;

        if (!data || !data.success || !data.result) {
            searchResults.value = [];
            showSearchResults.value = false;
            return;
        }

        const results: SearchResultItem[] = data.result.map((item: AssetInfoResponse) => ({
            id: item.id,
            code: item.code,
            name: item.name,
            category: item.category,
            currency: item.currency,
            market: item.market
        }));

        searchResults.value = results;
        showSearchResults.value = results.length > 0;

        // Load prices for search results
        loadSearchResultPrices(results);
    } catch (error) {
        logger.error('Failed to search assets', error);
        searchResults.value = [];
        showSearchResults.value = false;
    } finally {
        searchLoading.value = false;
    }
}

async function loadSearchResultPrices(items: SearchResultItem[]): Promise<void> {
    const pricePromises = items.map(async (item) => {
        try {
            const result = await investmentStore.loadLatestMarketData({ assetId: item.id });
            if (result) {
                item.currentPrice = result.price;
            }
        } catch {
            // Ignore individual failures
        }
    });

    await Promise.all(pricePromises);
}

function onSearchResultClick(item: SearchResultItem): void {
    showSearchResults.value = false;
    searchQuery.value = '';

    // Check if already in holdings or watchlist
    const existingHolding = holdingsList.value.find(h => h.assetId === item.id);
    if (existingHolding) {
        showDetail(existingHolding);
        return;
    }

    const existingWatchlist = watchlistAssets.value.find(w => w.assetId === item.id);
    if (existingWatchlist) {
        showDetail(existingWatchlist);
        return;
    }

    // Show as a temporary detail
    const displayItem: DisplayAsset = {
        assetId: item.id,
        assetCode: item.code,
        assetName: item.name,
        category: item.category,
        currency: item.currency,
        market: item.market,
        currentPrice: item.currentPrice,
        changeRate: item.changeRate,
        isHolding: false
    };
    showDetail(displayItem);
}

function onViewAllResults(): void {
    showSearchResults.value = false;
    // Keep the search query to filter the current tab
}

function isInWatchlist(assetId: string): boolean {
    return watchlistAssets.value.some(w => w.assetId === assetId);
}

async function addToWatchlist(item: SearchResultItem): Promise<void> {
    try {
        await investmentStore.addUserAsset({ assetId: item.id });
        // Reload watchlist
        await loadWatchlist();
    } catch (error) {
        logger.error('Failed to add to watchlist', error);
    }
}

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

// --- Action handlers ---
function onBuyClick(item: DisplayAsset | SearchResultItem): void {
    // TODO: Open buy transaction dialog
    logger.debug('Buy clicked for', item);
}

function onSellClick(item: DisplayAsset): void {
    // TODO: Open sell transaction dialog
    logger.debug('Sell clicked for', item);
}

function onManageClick(): void {
    // TODO: Open asset management dialog
    logger.debug('Manage clicked');
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

.page-body {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

/* Search Section */
.search-section {
    display: flex;
    align-items: center;
}

.search-wrapper {
    flex: 1;
    position: relative;
}

.search-input {
    max-width: 100%;
}

.manage-btn {
    flex-shrink: 0;
    height: 40px;
}

/* Search Results Overlay */
.search-results-overlay {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    z-index: 100;
    margin-top: 4px;
    max-height: 400px;
    overflow-y: auto;
}

/* Summary Footer */
.summary-item {
    display: flex;
    align-items: center;
    white-space: nowrap;
}

/* Tables */
.holdings-table :deep(.v-data-table__td),
.watchlist-table :deep(.v-data-table__td) {
    padding-top: 8px;
    padding-bottom: 8px;
}
</style>

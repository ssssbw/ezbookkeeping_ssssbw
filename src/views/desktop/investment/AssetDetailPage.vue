<template>
    <div class="asset-detail-page">
        <div class="page-header">
            <v-btn variant="text" :icon="mdiArrowLeft" @click="goBack" />
            <span class="text-h6">{{ displayName }}</span>
        </div>

        <div class="page-content">
            <!-- Basic Info -->
            <v-card class="mb-4">
                <v-card-text>
                    <v-row dense>
                        <v-col cols="6" md="3">
                            <div class="info-label">{{ tt('asset.AssetCode') }}</div>
                            <div class="info-value asset-code">{{ assetCode }}</div>
                        </v-col>
                        <v-col cols="6" md="3">
                            <div class="info-label">{{ tt('asset.AssetName') }}</div>
                            <div class="info-value">{{ assetName }}</div>
                        </v-col>
                        <v-col cols="6" md="3">
                            <div class="info-label">{{ tt('Market') }}</div>
                            <div class="info-value">{{ formatMarket(displayMarket, tt) }}</div>
                        </v-col>
                        <v-col cols="6" md="3">
                            <div class="info-label">{{ tt('Category') }}</div>
                            <div class="info-value">{{ formatCategory(displayCategory, tt) }}</div>
                        </v-col>
                        <v-col cols="6" md="3">
                            <div class="info-label">{{ tt('Current Price') }}</div>
                            <div class="info-value">{{ formatPriceWithDate(displayCurrentPrice, displayCurrentPriceDate) }}</div>
                        </v-col>
                        <v-col cols="6" md="3">
                            <div class="info-label">{{ tt('Total Quantity') }}</div>
                            <div class="info-value">{{ formatQuantity(displayTotalQuantity) }}</div>
                        </v-col>
                        <v-col cols="6" md="3">
                            <div class="info-label">{{ tt('Total Value') }}</div>
                            <div class="info-value">{{ formatCurrencyValue(displayTotalMarketValue, displayCurrency) }}</div>
                        </v-col>
                        <v-col cols="6" md="3">
                            <div class="info-label">{{ tt('Total Cost') }}</div>
                            <div class="info-value">{{ formatCurrencyValue(displayTotalCost, displayCurrency) }}</div>
                        </v-col>
                        <v-col cols="6" md="3">
                            <div class="info-label">{{ tt('Unrealized P&L') }}</div>
                            <div class="info-value" :class="getReturnColorClass(displayUnrealizedPnl)">
                                {{ formatCurrencyValue(displayUnrealizedPnl, displayCurrency) }}
                            </div>
                        </v-col>
                        <v-col cols="6" md="3">
                            <div class="info-label">{{ tt('Weighted Return') }}</div>
                            <div class="info-value" :class="getReturnColorClass(displayReturnRate)">
                                {{ formatReturnRate(displayReturnRate) }}
                            </div>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>

            <!-- Holding Detail -->
            <v-card class="mb-4" v-if="asset">
                <v-card-title class="text-subtitle-1">{{ tt('Holding Detail') }}</v-card-title>
                <v-card-text class="pa-0">
                    <v-data-table
                        :headers="holdingHeaders"
                        :items="asset.holdings"
                        :hover="true"
                        :items-per-page="-1"
                        class="holdings-detail-table"
                    >
                        <template #item.accountName="{ item }">
                            <span class="text-body-2">{{ item.accountName || item.accountId }}</span>
                        </template>
                        <template #item.currentPrice="{ item }">
                            <span class="text-body-2">{{ formatPriceWithDate(item.currentPrice, item.currentPriceDate) }}</span>
                        </template>
                        <template #item.quantity="{ item }">
                            <span class="text-body-2">{{ formatQuantity(item.quantity) }}</span>
                        </template>
                        <template #item.marketValue="{ item }">
                            <span class="text-body-2">{{ formatCurrencyValue(item.marketValue, item.currency) }}</span>
                        </template>
                        <template #item.totalCost="{ item }">
                            <span class="text-body-2">{{ formatCurrencyValue(item.totalCost, item.currency) }}</span>
                        </template>
                        <template #item.unrealizedPnl="{ item }">
                            <span class="text-body-2" :class="getReturnColorClass(item.unrealizedPnl)">
                                {{ formatCurrencyValue(item.unrealizedPnl, item.currency) }}
                            </span>
                        </template>
                        <template #item.returnRate="{ item }">
                            <span class="text-body-2" :class="getReturnColorClass(item.returnRate)">
                                {{ formatReturnRate(item.returnRate) }}
                            </span>
                        </template>
                        <template #no-data>
                            <div class="text-center py-4 text-medium-emphasis">
                                {{ tt('No holdings data') }}
                            </div>
                        </template>
                    </v-data-table>
                </v-card-text>
            </v-card>

            <!-- Price History -->
            <v-card class="mb-4">
                <v-card-title class="d-flex align-center justify-space-between">
                    <span class="text-subtitle-1">{{ tt('Price History') }}</span>
                    <div class="d-flex align-center ga-2">
                        <v-btn-toggle v-model="selectedTimeRange" density="compact" variant="outlined" divided>
                            <v-btn v-for="range in timeRanges" :key="range.value" :value="range.value" size="small">
                                {{ range.label }}
                            </v-btn>
                        </v-btn-toggle>
                        <v-btn size="small" color="primary" variant="tonal" :icon="mdiPlus" @click="openAddPrice" />
                    </div>
                </v-card-title>
                <v-card-text>
                    <v-chart v-if="priceChartData.length > 0" autoresize class="price-chart" :option="priceChartOptions" @click="onChartClick" />
                    <div v-else class="text-center py-6 text-medium-emphasis">
                        {{ tt('No price history data') }}
                    </div>
                </v-card-text>
            </v-card>

            <!-- Trade History -->
            <v-card>
                <v-card-title class="text-subtitle-1">{{ tt('Trade History') }}</v-card-title>
                <v-card-text class="pa-0">
                    <v-data-table
                        :headers="transactionHeaders"
                        :items="transactions"
                        :loading="transactionsLoading"
                        :hover="true"
                        :items-per-page="10"
                        class="transactions-table"
                    >
                        <template #item.tradeTime="{ item }">
                            <span class="text-body-2">{{ formatTradeTime(item.tradeTime) }}</span>
                        </template>
                        <template #item.type="{ item }">
                            <v-chip size="small" :color="getTransactionTypeColor(item.type)" variant="tonal">
                                {{ formatTransactionType(item.type, tt) }}
                            </v-chip>
                        </template>
                        <template #item.quantity="{ item }">
                            <span class="text-body-2">{{ formatQuantity(item.quantity) }}</span>
                        </template>
                        <template #item.price="{ item }">
                            <span class="text-body-2">{{ formatPrice(item.price) }}</span>
                        </template>
                        <template #item.amount="{ item }">
                            <span class="text-body-2">{{ formatCurrencyValue(item.amount, displayCurrency) }}</span>
                        </template>
                        <template #loading>
                            <v-skeleton-loader type="table-row@5" :loading="true" />
                        </template>
                        <template #no-data>
                            <div class="text-center py-6 text-medium-emphasis">
                                {{ tt('No transaction data') }}
                            </div>
                        </template>
                    </v-data-table>
                </v-card-text>
            </v-card>
        </div>

        <market-data-edit-dialog
            v-model="showPriceDialog"
            :asset-id="assetId"
            :currency="displayCurrency"
            :edit-date="editPriceDate"
            :edit-price="editPriceValue"
            :edit-volume="editPriceVolume"
            @saved="onPriceSaved"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { mdiArrowLeft, mdiPlus } from '@mdi/js';

import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';
import services from '@/lib/services.ts';

import { InvestmentTransactionType } from '@/models/investment.ts';
import MarketDataEditDialog from './components/MarketDataEditDialog.vue';

import {
    formatMarket, formatCategory, formatPrice, formatPriceWithDate, formatQuantity,
    formatCurrencyValue, formatReturnRate, getReturnColorClass
} from './assets/assetUtils.ts';

import type { AssetInfoResponse, InvestmentTransactionInfoResponse } from '@/models/investment.ts';

const PRIMARY_COLOR = '#c67e48';
const DIVISOR = 10000;

const route = useRoute();
const router = useRouter();
const { tt } = useI18n();
const investmentStore = useInvestmentStore();

const assetId = computed(() => route.params['id'] as string);

// Fallback data for watchlist-only assets
const fallbackAsset = ref<AssetInfoResponse | null>(null);
const fallbackPrice = ref<number | undefined>(undefined);
const fallbackPriceDate = ref<number | undefined>(undefined);

const asset = computed(() => {
    return investmentStore.aggregatedHoldings.find(h => h.assetId === assetId.value);
});

// Display values - from holdings if available, otherwise fallback
const assetCode = computed(() => asset.value?.assetCode || fallbackAsset.value?.code || '--');
const assetName = computed(() => asset.value?.assetName || fallbackAsset.value?.name || '--');
const displayMarket = computed(() => asset.value?.market ?? fallbackAsset.value?.market ?? 0);
const displayCategory = computed(() => asset.value?.category || fallbackAsset.value?.category || '');
const displayCurrency = computed(() => asset.value?.currency || fallbackAsset.value?.currency || 'CNY');
const displayName = computed(() => assetName.value !== '--' ? assetName.value : assetCode.value);

const displayCurrentPrice = computed(() => asset.value?.currentPrice ?? fallbackPrice.value);
const displayCurrentPriceDate = computed(() => asset.value?.currentPriceDate ?? fallbackPriceDate.value);
const displayTotalQuantity = computed(() => asset.value?.totalQuantity);
const displayTotalMarketValue = computed(() => asset.value?.totalMarketValue);
const displayTotalCost = computed(() => asset.value?.totalCost);
const displayUnrealizedPnl = computed(() => asset.value?.unrealizedPnl);
const displayReturnRate = computed(() => asset.value?.weightedReturnRate);

const transactions = ref<InvestmentTransactionInfoResponse[]>([]);
const transactionsLoading = ref(false);

const selectedTimeRange = ref('all');
const priceChartData = ref<{ date: string; price: number; isManual: boolean }[]>([]);

// Price edit dialog
const showPriceDialog = ref(false);
const editPriceDate = ref<string | undefined>(undefined);
const editPriceValue = ref<number | undefined>(undefined);
const editPriceVolume = ref<number | undefined>(undefined);

const timeRanges = computed(() => [
    { value: '1m', label: tt('1 Month') },
    { value: '3m', label: tt('3 Months') },
    { value: '6m', label: tt('6 Months') },
    { value: '1y', label: tt('1 Year') },
    { value: 'all', label: tt('All') },
]);

const holdingHeaders = computed(() => [
    { key: 'accountName', title: tt('Account'), sortable: false },
    { key: 'currentPrice', title: tt('Current Price'), sortable: false },
    { key: 'quantity', title: tt('Quantity'), sortable: false },
    { key: 'marketValue', title: tt('Market Value'), sortable: false },
    { key: 'totalCost', title: tt('Total Cost'), sortable: false },
    { key: 'unrealizedPnl', title: tt('Unrealized P&L'), sortable: false },
    { key: 'returnRate', title: tt('Return Rate'), sortable: false },
]);

const transactionHeaders = computed(() => [
    { key: 'tradeTime', title: tt('Date'), sortable: false },
    { key: 'type', title: tt('Type'), sortable: false },
    { key: 'quantity', title: tt('Quantity'), sortable: false },
    { key: 'price', title: tt('Price'), sortable: false },
    { key: 'amount', title: tt('Amount'), sortable: false },
]);

const priceChartOptions = computed(() => {
    const dates = priceChartData.value.map(d => d.date);
    const prices = priceChartData.value.map(d => d.price / DIVISOR);
    const manualIndices = priceChartData.value.map((d, i) => d.isManual ? i : -1).filter(i => i >= 0);

    return {
        grid: {
            left: 60,
            right: 20,
            top: 20,
            bottom: 30,
        },
        tooltip: {
            trigger: 'axis',
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            formatter: (params: any) => {
                const p = params[0];
                const idx = p.dataIndex;
                const item = priceChartData.value[idx];
                if (!item) return '';
                let tip = `${p.axisValue}<br/>${tt('Price')}: ${formatPrice(item.price)}`;
                if (item.isManual) {
                    tip += `<br/><span style="color:#ff9800">${tt('Manually Entered')}</span>`;
                }
                return tip;
            }
        },
        xAxis: {
            type: 'category',
            data: dates,
            axisLabel: { fontSize: 11 },
        },
        yAxis: {
            type: 'value',
            axisLabel: {
                fontSize: 11,
                formatter: (val: number) => val.toFixed(2),
            },
        },
        series: [
            {
                type: 'line',
                data: prices,
                smooth: true,
                symbol: 'circle',
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                symbolSize: (val: number, params: any) => {
                    return manualIndices.includes(params.dataIndex) ? 8 : 0;
                },
                itemStyle: { color: PRIMARY_COLOR },
                lineStyle: { width: 2 },
                areaStyle: {
                    color: {
                        type: 'linear',
                        x: 0, y: 0, x2: 0, y2: 1,
                        colorStops: [
                            { offset: 0, color: 'rgba(198, 126, 72, 0.3)' },
                            { offset: 1, color: 'rgba(198, 126, 72, 0.02)' },
                        ],
                    },
                },
                markPoint: {
                    data: manualIndices.map(i => ({
                        coord: [dates[i], prices[i]],
                        symbol: 'circle',
                        symbolSize: 8,
                        itemStyle: { color: '#ff9800' },
                    })),
                },
            },
        ],
    };
});

function formatTradeTime(timestamp: number): string {
    if (!timestamp) return '--';
    const d = new Date(timestamp * 1000);
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
}

function formatTransactionType(type: number, tt: (key: string) => string): string {
    switch (type) {
        case InvestmentTransactionType.Buy: return tt('Buy');
        case InvestmentTransactionType.Sell: return tt('Sell');
        case InvestmentTransactionType.DividendCash: return tt('Dividend');
        case InvestmentTransactionType.DividendReinvest: return tt('Dividend Reinvest');
        default: return type.toString();
    }
}

function getTransactionTypeColor(type: number): string {
    switch (type) {
        case InvestmentTransactionType.Buy: return 'primary';
        case InvestmentTransactionType.Sell: return 'error';
        case InvestmentTransactionType.DividendCash:
        case InvestmentTransactionType.DividendReinvest: return 'success';
        default: return 'default';
    }
}

function getTimeRangeStartTime(): number {
    const now = Math.floor(Date.now() / 1000);
    switch (selectedTimeRange.value) {
        case '1m': return now - 30 * 24 * 3600;
        case '3m': return now - 90 * 24 * 3600;
        case '6m': return now - 180 * 24 * 3600;
        case '1y': return now - 365 * 24 * 3600;
        default: return 0;
    }
}

async function loadPriceHistory(): Promise<void> {
    try {
        const startTime = getTimeRangeStartTime();
        const resp = await services.getMarketDataList({
            asset_id: assetId.value,
            start_time: startTime > 0 ? startTime : undefined,
        });
        const data = resp.data;
        if (data.success && data.result) {
            priceChartData.value = data.result.map((item: { date: number; price: number; isManual?: boolean }) => ({
                date: formatDate(item.date),
                price: item.price,
                isManual: item.isManual || false,
            }));
        }
    } catch {
        priceChartData.value = [];
    }
}

function formatDate(timestamp: number): string {
    const d = new Date(timestamp * 1000);
    return `${d.getMonth() + 1}/${d.getDate()}`;
}

async function loadTransactions(): Promise<void> {
    transactionsLoading.value = true;
    try {
        const resp = await services.getInvestmentTransactions({
            asset_id: assetId.value,
        });
        const data = resp.data;
        if (data.success && data.result) {
            transactions.value = data.result;
        }
    } catch {
        transactions.value = [];
    } finally {
        transactionsLoading.value = false;
    }
}

async function loadFallbackAsset(): Promise<void> {
    try {
        const [assetResp, priceResp] = await Promise.all([
            services.getGlobalAsset({ id: assetId.value }),
            services.getMarketDataList({ asset_id: assetId.value }),
        ]);
        const assetData = assetResp.data;
        if (assetData.success && assetData.result) {
            fallbackAsset.value = assetData.result;
        }
        const priceData = priceResp.data;
        if (priceData.success && priceData.result && priceData.result.length > 0) {
            const latest = priceData.result[priceData.result.length - 1];
            if (latest) {
                fallbackPrice.value = latest.price;
                fallbackPriceDate.value = latest.date;
            }
        }
    } catch {
        // ignore
    }
}

function openAddPrice(): void {
    editPriceDate.value = undefined;
    editPriceValue.value = undefined;
    editPriceVolume.value = undefined;
    showPriceDialog.value = true;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function onChartClick(params: any): void {
    if (params.componentType !== 'series') return;
    const idx = params.dataIndex;
    const item = priceChartData.value[idx];
    if (!item) return;

    const parts = item.date.split('/');
    const month = parts[0];
    const day = parts[1];
    if (!month || !day) return;
    const year = new Date().getFullYear();
    editPriceDate.value = `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`;
    editPriceValue.value = item.price;
    editPriceVolume.value = undefined;
    showPriceDialog.value = true;
}

function onPriceSaved(): void {
    loadPriceHistory();
}

function goBack(): void {
    router.push('/investment/assets');
}

watch(selectedTimeRange, () => {
    loadPriceHistory();
});

onMounted(async () => {
    if (!asset.value) {
        await loadFallbackAsset();
    }
    loadPriceHistory();
    loadTransactions();
});
</script>

<style scoped>
.asset-detail-page {
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow-y: auto;
}

.page-header {
    padding: 16px 24px;
}

.page-content {
    padding: 0 24px 24px;
}

.info-label {
    font-size: 12px;
    color: rgba(var(--v-theme-on-surface), 0.6);
    margin-bottom: 2px;
}

.info-value {
    font-size: 14px;
}

.asset-code {
    color: rgb(var(--v-theme-primary));
    font-weight: 500;
}

.price-chart {
    width: 100%;
    height: 300px;
}

.holdings-detail-table :deep(.v-data-table__td) {
    padding-top: 8px;
    padding-bottom: 8px;
}

.transactions-table :deep(.v-data-table__td) {
    padding-top: 8px;
    padding-bottom: 8px;
}
</style>

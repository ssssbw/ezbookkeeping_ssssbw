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
                    <v-divider class="mt-3" />
                    <div class="mt-3">
                        <v-btn variant="text" color="primary" size="small" @click="goToTransactions">
                            {{ tt('Trade History') }}
                        </v-btn>
                    </div>
                </v-card-text>
            </v-card>

            <!-- Price History -->
            <v-card class="mb-4">
                <v-card-title class="d-flex align-center justify-space-between">
                    <span class="text-subtitle-1">{{ tt('Price History') }}</span>
                    <div class="d-flex align-center ga-2">
                        <v-btn-toggle v-model="selectedTimeRange" density="compact" variant="outlined" divided>
                            <v-btn v-for="range in timeRanges" :key="range.value" :value="range.value" size="small" class="px-3">
                                {{ range.label }}
                            </v-btn>
                        </v-btn-toggle>
                        <v-btn size="small" color="primary" variant="tonal" @click="openAddPrice">
                            {{ tt('Add') }}
                        </v-btn>
                        <v-btn size="small" variant="tonal" @click="refreshPriceHistory">
                            {{ tt('Refresh') }}
                        </v-btn>
                    </div>
                </v-card-title>
                <v-card-text>
                    <v-chart v-if="priceChartData.length > 0" autoresize class="price-chart" :option="priceChartOptions" @click="onChartClick" />
                    <div v-else class="text-center py-6 text-medium-emphasis">
                        {{ tt('No price history data') }}
                    </div>
                </v-card-text>
                <v-divider />
                <v-card-text class="pa-0">
                    <v-data-table
                        :headers="priceListHeaders"
                        :items="priceListData"
                        :hover="true"
                        items-per-page="10"
                        class="price-list-table"
                    >
                        <template #item.date="{ item }">
                            <span class="text-body-2">{{ item.date }}</span>
                        </template>
                        <template #item.price="{ item }">
                            <span class="text-body-2">{{ formatPrice(item.price) }}</span>
                        </template>
                        <template #item.volume="{ item }">
                            <span class="text-body-2">{{ item.volume ? formatQuantity(item.volume) : '--' }}</span>
                        </template>
                        <template #item.isManual="{ item }">
                            <v-chip v-if="item.isManual" size="small" color="warning" variant="tonal">
                                {{ tt('Manual') }}
                            </v-chip>
                        </template>
                        <template #item.actions="{ item }">
                            <v-btn size="small" variant="text" color="primary" @click.stop="editPrice(item)">
                                {{ tt('Edit') }}
                            </v-btn>
                        </template>
                        <template #no-data>
                            <div class="text-center py-4 text-medium-emphasis">
                                {{ tt('No price history data') }}
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
import { mdiArrowLeft } from '@mdi/js';

import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';
import services from '@/lib/services.ts';

import MarketDataEditDialog from './components/MarketDataEditDialog.vue';

import {
    formatMarket, formatCategory, formatPrice, formatPriceWithDate, formatQuantity,
    formatCurrencyValue, formatReturnRate, getReturnColorClass
} from './assets/assetUtils.ts';

import type { AssetInfoResponse } from '@/models/investment.ts';

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

const selectedTimeRange = ref('all');
const priceChartData = ref<{ date: string; price: number; isManual: boolean }[]>([]);
const priceListRawData = ref<{ date: number; price: number; volume: number; isManual: boolean }[]>([]);

// Price list table
const priceListHeaders = computed(() => [
    { key: 'date', title: tt('Date'), sortable: false },
    { key: 'price', title: tt('Price'), sortable: false },
    { key: 'volume', title: tt('Volume'), sortable: false },
    { key: 'isManual', title: '', sortable: false, width: '80' },
    { key: 'actions', title: '', sortable: false, width: '80' },
]);

const priceListData = computed(() => {
    return priceListRawData.value.map(item => ({
        ...item,
        date: formatDateFull(item.date),
    }));
});

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
            priceListRawData.value = data.result.map((item: { date: number; price: number; volume?: number; isManual?: boolean }) => ({
                date: item.date,
                price: item.price,
                volume: item.volume || 0,
                isManual: item.isManual || false,
            }));
            priceChartData.value = priceListRawData.value.map(item => ({
                date: formatDate(item.date),
                price: item.price,
                isManual: item.isManual,
            }));
        }
    } catch {
        priceChartData.value = [];
        priceListRawData.value = [];
    }
}

function formatDate(timestamp: number): string {
    const d = new Date(timestamp * 1000);
    return `${d.getMonth() + 1}/${d.getDate()}`;
}

function formatDateFull(timestamp: number): string {
    const d = new Date(timestamp * 1000);
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
}

function editPrice(item: { date: string; price: number; volume?: number }): void {
    editPriceDate.value = item.date;
    editPriceValue.value = item.price;
    editPriceVolume.value = item.volume;
    showPriceDialog.value = true;
}

async function refreshPriceHistory(): Promise<void> {
    await loadPriceHistory();
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

function goToTransactions(): void {
    router.push({ path: '/investment/transactions', query: { assetId: assetId.value } });
}

watch(selectedTimeRange, () => {
    loadPriceHistory();
});

onMounted(async () => {
    if (!asset.value) {
        await loadFallbackAsset();
    }
    loadPriceHistory();
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

.price-list-table :deep(.v-data-table__td) {
    padding-top: 8px;
    padding-bottom: 8px;
}
</style>

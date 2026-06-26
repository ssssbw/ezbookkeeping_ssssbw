<template>
    <div class="page-content">
        <div class="page-header">
            <h1 class="page-title">{{ tt('Investment Portfolio') }}</h1>
        </div>
        <div class="page-body">
            <v-row>
                <v-col cols="12" md="3">
                    <v-card>
                        <v-card-text class="summary-card">
                            <div class="summary-label">{{ tt('Current Value') }}</div>
                            <div class="summary-value">{{ formatCurrencyValue(holdingsSummary.totalMarketValue, primaryCurrency) }}</div>
                        </v-card-text>
                    </v-card>
                </v-col>
                <v-col cols="12" md="3">
                    <v-card>
                        <v-card-text class="summary-card">
                            <div class="summary-label">{{ tt('Total Cost') }}</div>
                            <div class="summary-value">{{ formatCurrencyValue(holdingsSummary.totalCost, primaryCurrency) }}</div>
                        </v-card-text>
                    </v-card>
                </v-col>
                <v-col cols="12" md="3">
                    <v-card>
                        <v-card-text class="summary-card">
                            <div class="summary-label">{{ tt('Unrealized P&L') }}</div>
                            <div class="summary-value" :class="getReturnColorClass(holdingsSummary.totalUnrealizedPnl)">
                                {{ formatCurrencyValue(holdingsSummary.totalUnrealizedPnl, primaryCurrency) }}
                            </div>
                        </v-card-text>
                    </v-card>
                </v-col>
                <v-col cols="12" md="3">
                    <v-card>
                        <v-card-text class="summary-card">
                            <div class="summary-label">{{ tt('Return Rate') }}</div>
                            <div class="summary-value" :class="getReturnColorClass(holdingsSummary.totalReturnRate)">
                                {{ formatReturnRate(holdingsSummary.totalReturnRate) }}
                            </div>
                        </v-card-text>
                    </v-card>
                </v-col>
            </v-row>

            <v-row>
                <v-col cols="12" lg="4" md="5">
                    <v-card>
                        <v-card-title class="text-subtitle-1">{{ tt('Asset Allocation') }}</v-card-title>
                        <v-card-text>
                            <v-chart v-if="allocationData.length > 0" autoresize class="allocation-chart" :option="allocationChartOptions" />
                            <div v-else class="text-center py-6 text-medium-emphasis">
                                {{ tt('No holdings data') }}
                            </div>
                            <div v-if="allocationData.length > 0" class="mt-4">
                                <div v-for="(item, index) in allocationData" :key="index" class="d-flex align-center mb-2">
                                    <v-icon :color="item.color" size="12" class="me-2" :icon="mdiCheckboxMarkedCircle" />
                                    <span class="text-body-2 me-auto">{{ item.name }}</span>
                                    <span class="text-body-2">{{ item.percent }}</span>
                                </div>
                            </div>
                        </v-card-text>
                    </v-card>
                </v-col>

                <v-col cols="12" lg="8" md="7">
                    <v-card>
                        <v-card-title class="text-subtitle-1">{{ tt('Holdings') }}</v-card-title>
                        <v-card-text class="pa-0">
                            <v-data-table
                                :headers="headers"
                                :items="aggregatedHoldings"
                                :hover="true"
                                :items-per-page="-1"
                                class="holdings-table"
                            >
                                <template #item.assetCode="{ item }">
                                    <span class="text-body-2 asset-code">{{ item.assetCode }}</span>
                                </template>
                                <template #item.assetName="{ item }">
                                    <span class="text-body-2">{{ item.assetName }}</span>
                                </template>
                                <template #item.category="{ item }">
                                    <v-chip size="small" :color="getCategoryColor(item.category)" variant="tonal">
                                        {{ formatCategory(item.category, tt) }}
                                    </v-chip>
                                </template>
                                <template #item.totalQuantity="{ item }">
                                    <span class="text-body-2">{{ formatQuantity(item.totalQuantity) }}</span>
                                </template>
                                <template #item.currentPrice="{ item }">
                                    <span class="text-body-2">{{ formatPriceWithDate(item.currentPrice, item.currentPriceDate) }}</span>
                                </template>
                                <template #item.totalMarketValue="{ item }">
                                    <span class="text-body-2">{{ formatCurrencyValue(item.totalMarketValue, item.currency) }}</span>
                                </template>
                                <template #item.totalCost="{ item }">
                                    <span class="text-body-2">{{ formatCurrencyValue(item.totalCost, item.currency) }}</span>
                                </template>
                                <template #item.unrealizedPnl="{ item }">
                                    <span class="text-body-2" :class="getReturnColorClass(item.unrealizedPnl)">
                                        {{ formatCurrencyValue(item.unrealizedPnl, item.currency) }}
                                    </span>
                                </template>
                                <template #item.weightedReturnRate="{ item }">
                                    <span class="text-body-2" :class="getReturnColorClass(item.weightedReturnRate)">
                                        {{ formatReturnRate(item.weightedReturnRate) }}
                                    </span>
                                </template>
                                <template #no-data>
                                    <div class="text-center py-6 text-medium-emphasis">
                                        {{ tt('No holdings data') }}
                                    </div>
                                </template>
                            </v-data-table>
                        </v-card-text>
                    </v-card>
                </v-col>
            </v-row>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';

import { mdiCheckboxMarkedCircle } from '@mdi/js';

import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';

import {
    formatPriceWithDate, formatQuantity, formatCurrencyValue,
    formatReturnRate, formatCategory, getReturnColorClass, getCategoryColor
} from './assets/assetUtils.ts';

const { tt } = useI18n();
const investmentStore = useInvestmentStore();

const loading = ref(true);

const aggregatedHoldings = computed(() => investmentStore.aggregatedHoldings);
const holdingsSummary = computed(() => investmentStore.holdingsSummary);

const primaryCurrency = computed(() => {
    const first = aggregatedHoldings.value[0];
    return first ? first.currency : 'CNY';
});

interface AllocationItem {
    name: string;
    value: number;
    percent: string;
    color: string;
}

const allocationData = computed<AllocationItem[]>(() => {
    const categoryMap = new Map<string, number>();
    for (const h of aggregatedHoldings.value) {
        const cat = h.category || 'other';
        categoryMap.set(cat, (categoryMap.get(cat) || 0) + h.totalMarketValue);
    }

    const total = Array.from(categoryMap.values()).reduce((a, b) => a + b, 0);
    if (total === 0) return [];

    const colorMap: Record<string, string> = {
        equity: 'primary',
        fixed_income: 'success',
        commodity: 'warning',
        digital: 'purple',
        other: 'grey',
    };

    return Array.from(categoryMap.entries())
        .map(([cat, value]) => ({
            name: formatCategory(cat, tt),
            value,
            percent: ((value / total) * 100).toFixed(1) + '%',
            color: colorMap[cat] || 'grey',
        }))
        .sort((a, b) => b.value - a.value);
});

const allocationChartOptions = computed(() => {
    const PIE_COLORS = ['#c67e48', '#4caf50', '#ff9800', '#9c27b0', '#607d8b'];
    return {
        tooltip: {
            trigger: 'item',
            formatter: '{b}: {d}%',
        },
        series: [
            {
                type: 'pie',
                radius: ['40%', '70%'],
                avoidLabelOverlap: true,
                itemStyle: {
                    borderRadius: 4,
                    borderColor: '#fff',
                    borderWidth: 2,
                },
                label: {
                    show: false,
                },
                emphasis: {
                    label: {
                        show: true,
                        fontSize: 14,
                        fontWeight: 'bold',
                    },
                },
                data: allocationData.value.map((item, i) => ({
                    value: item.value,
                    name: item.name,
                    itemStyle: { color: PIE_COLORS[i % PIE_COLORS.length] },
                })),
            },
        ],
    };
});

const headers = computed(() => [
    { key: 'assetCode', title: tt('asset.AssetCode'), sortable: false, width: '100' },
    { key: 'assetName', title: tt('asset.AssetName'), sortable: false },
    { key: 'category', title: tt('Category'), sortable: false, width: '100' },
    { key: 'totalQuantity', title: tt('Quantity'), sortable: false },
    { key: 'currentPrice', title: tt('Current Price'), sortable: false },
    { key: 'totalMarketValue', title: tt('Market Value'), sortable: false },
    { key: 'totalCost', title: tt('Total Cost'), sortable: false },
    { key: 'unrealizedPnl', title: tt('Unrealized P&L'), sortable: false },
    { key: 'weightedReturnRate', title: tt('Return Rate'), sortable: false },
]);

onMounted(async () => {
    loading.value = true;
    try {
        await investmentStore.loadHoldings({ force: false });
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

.summary-card {
    text-align: center;
}

.summary-label {
    font-size: 12px;
    color: rgba(var(--v-theme-on-surface), 0.6);
    margin-bottom: 4px;
}

.summary-value {
    font-size: 20px;
    font-weight: 500;
}

.allocation-chart {
    width: 100%;
    height: 240px;
}

.holdings-table :deep(.v-data-table__td) {
    padding-top: 8px;
    padding-bottom: 8px;
}

.asset-code {
    color: rgb(var(--v-theme-primary));
    font-weight: 500;
}
</style>

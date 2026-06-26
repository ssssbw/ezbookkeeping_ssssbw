<template>
    <div class="page-content">
        <div class="page-header">
            <h1 class="page-title">{{ tt('Performance Analysis') }}</h1>
        </div>
        <div class="page-body">
            <v-row>
                <v-col cols="12" md="4">
                    <v-card>
                        <v-card-text class="metric-card">
                            <div class="metric-label">{{ tt('Return Rate') }}</div>
                            <div class="metric-value" :class="getReturnColorClass(overview?.totalReturnRate || 0)">
                                {{ formatReturnRate(overview?.totalReturnRate || 0) }}
                            </div>
                        </v-card-text>
                    </v-card>
                </v-col>
                <v-col cols="12" md="4">
                    <v-card>
                        <v-card-text class="metric-card">
                            <div class="metric-label">{{ tt('Unrealized P&L') }}</div>
                            <div class="metric-value" :class="getReturnColorClass(overview?.totalUnrealizedPnl || 0)">
                                {{ formatCurrencyValue(overview?.totalUnrealizedPnl || 0, 'CNY') }}
                            </div>
                        </v-card-text>
                    </v-card>
                </v-col>
                <v-col cols="12" md="4">
                    <v-card>
                        <v-card-text class="metric-card">
                            <div class="metric-label">{{ tt('Holdings Count') }}</div>
                            <div class="metric-value">{{ holdingsCount }}</div>
                        </v-card-text>
                    </v-card>
                </v-col>
            </v-row>

            <v-row>
                <v-col cols="12" md="6">
                    <v-card>
                        <v-card-title class="text-subtitle-1">{{ tt('Asset Class Performance') }}</v-card-title>
                        <v-card-text>
                            <v-chart v-if="categoryData.length > 0" autoresize class="analysis-chart" :option="categoryChartOptions" />
                            <div v-else class="text-center py-6 text-medium-emphasis">
                                {{ tt('No holdings data') }}
                            </div>
                        </v-card-text>
                    </v-card>
                </v-col>

                <v-col cols="12" md="6">
                    <v-card>
                        <v-card-title class="text-subtitle-1">{{ tt('Holdings Return') }}</v-card-title>
                        <v-card-text>
                            <v-chart v-if="holdingReturnData.length > 0" autoresize class="analysis-chart" :option="holdingReturnChartOptions" />
                            <div v-else class="text-center py-6 text-medium-emphasis">
                                {{ tt('No holdings data') }}
                            </div>
                        </v-card-text>
                    </v-card>
                </v-col>
            </v-row>

            <v-card>
                <v-card-title class="text-subtitle-1">{{ tt('Performance Metrics') }}</v-card-title>
                <v-card-text>
                    <v-row>
                        <v-col cols="12" md="3">
                            <div class="metric-detail">
                                <div class="metric-detail-label">{{ tt('Total Investment') }}</div>
                                <div class="metric-detail-value">{{ formatCurrencyValue(overview?.totalInvestment || 0, 'CNY') }}</div>
                            </div>
                        </v-col>
                        <v-col cols="12" md="3">
                            <div class="metric-detail">
                                <div class="metric-detail-label">{{ tt('Current Value') }}</div>
                                <div class="metric-detail-value">{{ formatCurrencyValue(overview?.totalMarketValue || 0, 'CNY') }}</div>
                            </div>
                        </v-col>
                        <v-col cols="12" md="3">
                            <div class="metric-detail">
                                <div class="metric-detail-label">{{ tt('Total Cost') }}</div>
                                <div class="metric-detail-value">{{ formatCurrencyValue(holdingsSummary.totalCost, 'CNY') }}</div>
                            </div>
                        </v-col>
                        <v-col cols="12" md="3">
                            <div class="metric-detail">
                                <div class="metric-detail-label">{{ tt('Weighted Return') }}</div>
                                <div class="metric-detail-value" :class="getReturnColorClass(holdingsSummary.totalReturnRate)">
                                    {{ formatReturnRate(holdingsSummary.totalReturnRate) }}
                                </div>
                            </div>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';

import {
    formatCurrencyValue, formatReturnRate, getReturnColorClass
} from './assets/assetUtils.ts';

const { tt } = useI18n();
const investmentStore = useInvestmentStore();

const loading = ref(true);

const overview = computed(() => investmentStore.overview);
const aggregatedHoldings = computed(() => investmentStore.aggregatedHoldings);
const holdingsSummary = computed(() => investmentStore.holdingsSummary);
const holdingsCount = computed(() => aggregatedHoldings.value.length);

interface CategoryData {
    name: string;
    value: number;
    returnRate: number;
}

const categoryData = computed<CategoryData[]>(() => {
    const map = new Map<string, { value: number; returnRate: number }>();
    for (const h of aggregatedHoldings.value) {
        const cat = h.category || 'other';
        const existing = map.get(cat);
        if (existing) {
            existing.value += h.totalMarketValue;
        } else {
            map.set(cat, { value: h.totalMarketValue, returnRate: h.weightedReturnRate });
        }
    }

    return Array.from(map.entries()).map(([cat, data]) => ({
        name: formatCategoryName(cat),
        value: data.value,
        returnRate: data.returnRate,
    })).sort((a, b) => b.value - a.value);
});

function formatCategoryName(category: string): string {
    switch (category) {
        case 'equity': return tt('Equity');
        case 'fixed_income': return tt('Fixed Income');
        case 'commodity': return tt('Commodity');
        case 'digital': return tt('Digital');
        default: return category;
    }
}

const holdingReturnData = computed(() => {
    return aggregatedHoldings.value
        .map(h => ({
            name: h.assetCode,
            returnRate: h.weightedReturnRate,
        }))
        .sort((a, b) => b.returnRate - a.returnRate);
});

const categoryChartOptions = computed(() => {
    const PIE_COLORS = ['#c67e48', '#4caf50', '#ff9800', '#9c27b0', '#607d8b'];
    return {
        tooltip: {
            trigger: 'item',
            formatter: '{b}: {c} ({d}%)'
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
                label: { show: false },
                emphasis: {
                    label: { show: true, fontSize: 14, fontWeight: 'bold' },
                },
                data: categoryData.value.map((item, i) => ({
                    value: item.value,
                    name: item.name,
                    itemStyle: { color: PIE_COLORS[i % PIE_COLORS.length] },
                })),
            },
        ],
    };
});

const holdingReturnChartOptions = computed(() => {
    const names = holdingReturnData.value.map(d => d.name);
    const values = holdingReturnData.value.map(d => d.returnRate / 100);

    return {
        tooltip: {
            trigger: 'axis',
            axisPointer: { type: 'shadow' },
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            formatter: (params: any) => {
                const p = Array.isArray(params) ? params[0] : params;
                if (!p) return '';
                return `${p.name}: ${p.value >= 0 ? '+' : ''}${Number(p.value).toFixed(2)}%`;
            }
        },
        grid: {
            left: '3%', right: '4%', bottom: '3%', top: '10%',
            containLabel: true,
        },
        xAxis: {
            type: 'category',
            data: names,
            axisLabel: { rotate: names.length > 5 ? 45 : 0 },
        },
        yAxis: {
            type: 'value',
            axisLabel: {
                formatter: (val: number) => val.toFixed(0) + '%',
            },
            splitLine: { lineStyle: { type: 'dashed' } },
        },
        series: [
            {
                type: 'bar',
                data: values.map(v => ({
                    value: v,
                    itemStyle: {
                        color: v >= 0
                            ? 'rgba(198, 126, 72, 0.8)'
                            : 'rgba(76, 175, 80, 0.8)',
                    },
                })),
                barWidth: '50%',
            },
        ],
    };
});

onMounted(async () => {
    loading.value = true;
    try {
        await Promise.all([
            investmentStore.loadOverview({ force: false }),
            investmentStore.loadHoldings({ force: false }),
        ]);
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

.metric-card {
    text-align: center;
}

.metric-label {
    font-size: 12px;
    color: rgba(var(--v-theme-on-surface), 0.6);
    margin-bottom: 4px;
}

.metric-value {
    font-size: 24px;
    font-weight: 500;
}

.analysis-chart {
    width: 100%;
    height: 300px;
}

.metric-detail {
    padding: 12px 0;
}

.metric-detail-label {
    font-size: 12px;
    color: rgba(var(--v-theme-on-surface), 0.6);
    margin-bottom: 4px;
}

.metric-detail-value {
    font-size: 16px;
    font-weight: 500;
}
</style>

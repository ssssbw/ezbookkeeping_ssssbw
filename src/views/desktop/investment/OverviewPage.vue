<template>
    <div class="page-content">
        <div class="page-header">
            <div class="d-flex align-center">
                <h1 class="page-title">{{ tt('Investment Overview') }}</h1>
                <v-btn density="compact" color="default" variant="text" size="24"
                       class="ms-2" :icon="true" :loading="loading" @click="reload(true)">
                    <template #loader>
                        <v-progress-circular indeterminate size="20"/>
                    </template>
                    <v-icon :icon="mdiRefresh" size="24" />
                    <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                </v-btn>
            </div>
        </div>
        <div class="page-body">
            <v-row class="match-height">
                <v-col cols="12" lg="6" md="12">
                    <v-card :class="{ 'disabled': loading }">
                        <template #title>
                            <span class="text-xl font-weight-bold">{{ tt('Asset Allocation') }}</span>
                        </template>
                        <v-card-text>
                            <div class="d-flex flex-column align-center">
                                <v-chart v-if="allocationData.length > 0" autoresize class="asset-allocation-chart" :option="assetAllocationChartOptions" />
                                <div v-else class="text-center py-6 text-medium-emphasis">
                                    {{ tt('No holdings data') }}
                                </div>
                            </div>
                            <div v-if="allocationData.length > 0" class="mt-4">
                                <div class="d-flex align-center mb-2" v-for="(item, index) in allocationData" :key="index">
                                    <v-icon :color="item.color" size="12" class="me-2" :icon="mdiCheckboxMarkedCircle" />
                                    <span class="text-body-2 me-2">{{ item.name }}</span>
                                    <span class="text-body-2 font-weight-bold">{{ item.percent }}</span>
                                </div>
                            </div>
                        </v-card-text>
                    </v-card>
                </v-col>

                <v-col cols="12" lg="6" md="12">
                    <v-card :class="{ 'disabled': loading }">
                        <template #title>
                            <span class="text-xl font-weight-bold">{{ tt('Portfolio Summary') }}</span>
                        </template>
                        <v-card-text>
                            <v-row>
                                <v-col cols="12" md="6">
                                    <div class="d-flex flex-column">
                                        <span class="text-caption text-medium-emphasis">{{ tt('Total Investment') }}</span>
                                        <span class="text-h5 font-weight-bold" v-if="!loading">{{ formatCurrencyValue(overview?.totalInvestment || 0, 'CNY') }}</span>
                                        <v-skeleton-loader class="mt-1" width="100px" type="text" :loading="true" v-else />
                                    </div>
                                </v-col>
                                <v-col cols="12" md="6">
                                    <div class="d-flex flex-column">
                                        <span class="text-caption text-medium-emphasis">{{ tt('Current Value') }}</span>
                                        <span class="text-h5 font-weight-bold" v-if="!loading">{{ formatCurrencyValue(overview?.totalMarketValue || 0, 'CNY') }}</span>
                                        <v-skeleton-loader class="mt-1" width="100px" type="text" :loading="true" v-else />
                                    </div>
                                </v-col>
                            </v-row>

                            <v-divider class="my-4" />

                            <v-row>
                                <v-col cols="12" md="6">
                                    <div class="d-flex flex-column">
                                        <span class="text-caption text-medium-emphasis">{{ tt('Return Rate') }}</span>
                                        <span class="text-h6 font-weight-bold" :class="getReturnColorClass(overview?.totalReturnRate || 0)" v-if="!loading">
                                            {{ formatReturnRate(overview?.totalReturnRate || 0) }}
                                        </span>
                                        <v-skeleton-loader class="mt-1" width="80px" type="text" :loading="true" v-else />
                                    </div>
                                </v-col>
                                <v-col cols="12" md="6">
                                    <div class="d-flex flex-column">
                                        <span class="text-caption text-medium-emphasis">{{ tt('Unrealized P&L') }}</span>
                                        <span class="text-h6 font-weight-bold" :class="getReturnColorClass(overview?.totalUnrealizedPnl || 0)" v-if="!loading">
                                            {{ formatCurrencyValue(overview?.totalUnrealizedPnl || 0, 'CNY') }}
                                        </span>
                                        <v-skeleton-loader class="mt-1" width="80px" type="text" :loading="true" v-else />
                                    </div>
                                </v-col>
                            </v-row>

                            <v-divider class="my-4" />

                            <div class="d-flex align-center justify-space-between">
                                <div class="d-flex flex-column">
                                    <span class="text-caption text-medium-emphasis">{{ tt('Holdings Count') }}</span>
                                    <span class="text-h5 font-weight-bold" v-if="!loading">{{ holdingsCount }}</span>
                                    <v-skeleton-loader class="mt-1" width="40px" type="text" :loading="true" v-else />
                                </div>
                            </div>
                        </v-card-text>
                    </v-card>
                </v-col>

                <v-col cols="12" md="6">
                    <v-row>
                        <v-col cols="6">
                            <investment-return-overview-card
                                :loading="loading" :disabled="loading" :icon="mdiCalendarTodayOutline"
                                :title="tt('Today')"
                                :return-amount="todayReturnDisplay"
                                :datetime="todayDate"
                            />
                        </v-col>
                        <v-col cols="6">
                            <investment-return-overview-card
                                :loading="loading" :disabled="loading" :icon="mdiCalendarWeekOutline"
                                :title="tt('This Week')"
                                :return-amount="'--'"
                                :datetime="weekDateRange"
                            />
                        </v-col>
                        <v-col cols="6">
                            <investment-return-overview-card
                                :loading="loading" :disabled="loading" :icon="mdiCalendarMonthOutline"
                                :title="tt('This Month')"
                                :return-amount="'--'"
                                :datetime="monthDateRange"
                            />
                        </v-col>
                        <v-col cols="6">
                            <investment-return-overview-card
                                :loading="loading" :disabled="loading" :icon="mdiLayersTripleOutline"
                                :title="tt('This Year')"
                                :return-amount="'--'"
                                :datetime="String(currentYear)"
                            />
                        </v-col>
                    </v-row>
                </v-col>

                <v-col cols="12" md="6">
                    <v-card :class="{ 'disabled': loading }">
                        <template #title>
                            <div class="d-flex align-center justify-space-between w-100">
                                <span class="text-xl font-weight-bold">{{ tt('Monthly Performance') }}</span>
                            </div>
                        </template>
                        <v-card-text>
                            <v-chart v-if="monthlyData.length > 0" autoresize class="monthly-performance-chart" :option="monthlyPerformanceChartOptions" />
                            <div v-else class="text-center py-6 text-medium-emphasis">
                                {{ tt('No holdings data') }}
                            </div>
                        </v-card-text>
                    </v-card>
                </v-col>
            </v-row>

            <snack-bar ref="snackbar" />
        </div>
    </div>
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';
import InvestmentReturnOverviewCard from '@/views/desktop/investment/components/InvestmentReturnOverviewCard.vue';

import { ref, computed, useTemplateRef, onMounted } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';

import { ThemeType } from '@/core/theme.ts';
import { formatCurrencyValue, formatReturnRate, getReturnColorClass } from './assets/assetUtils.ts';

import * as echarts from 'echarts/core';

import {
    mdiRefresh,
    mdiCheckboxMarkedCircle,
    mdiCalendarTodayOutline,
    mdiCalendarWeekOutline,
    mdiCalendarMonthOutline,
    mdiLayersTripleOutline
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const theme = useTheme();
const { tt } = useI18n();
const investmentStore = useInvestmentStore();
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const loading = ref<boolean>(true);

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

const overview = computed(() => investmentStore.overview);
const aggregatedHoldings = computed(() => investmentStore.aggregatedHoldings);
const holdingsCount = computed(() => aggregatedHoldings.value.length);

const todayDate = computed(() => {
    const d = new Date();
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
});

const currentYear = computed(() => new Date().getFullYear());

const weekDateRange = computed(() => {
    const now = new Date();
    const dayOfWeek = now.getDay();
    const monday = new Date(now);
    monday.setDate(now.getDate() - (dayOfWeek === 0 ? 6 : dayOfWeek - 1));
    const sunday = new Date(monday);
    sunday.setDate(monday.getDate() + 6);
    return `${monday.getMonth() + 1}/${monday.getDate()}-${sunday.getMonth() + 1}/${sunday.getDate()}`;
});

const monthDateRange = computed(() => {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
});

const todayReturnDisplay = computed(() => '--');

interface AllocationDisplayItem {
    name: string;
    value: number;
    percent: string;
    color: string;
}

const allocationData = computed<AllocationDisplayItem[]>(() => {
    const allocs = overview.value?.allocations;
    if (!allocs || allocs.length === 0) return [];

    const categoryColorMap: Record<string, string> = {
        equity: '#c67e48',
        fixed_income: '#4caf50',
        commodity: '#ff9800',
        digital: '#9c27b0',
        other: '#607d8b',
    };

    return allocs.map(a => ({
        name: formatCategoryName(a.category),
        value: a.value,
        percent: (a.percentage / 100).toFixed(1) + '%',
        color: categoryColorMap[a.category] || '#607d8b',
    }));
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

const assetAllocationChartOptions = computed(() => ({
    tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
    },
    legend: { show: false },
    series: [
        {
            type: 'pie',
            radius: ['45%', '75%'],
            avoidLabelOverlap: false,
            itemStyle: {
                borderRadius: 8,
                borderColor: isDarkMode.value ? '#1e1e1e' : '#fff',
                borderWidth: 2
            },
            label: { show: false, position: 'center' },
            emphasis: {
                label: { show: true, fontSize: 18, fontWeight: 'bold' }
            },
            labelLine: { show: false },
            data: allocationData.value.map(item => ({
                value: item.value,
                name: item.name,
                itemStyle: { color: item.color }
            }))
        }
    ]
}));

const monthlyData = computed(() => {
    if (aggregatedHoldings.value.length === 0) return [];
    return aggregatedHoldings.value.map(h => ({
        name: h.assetCode,
        value: h.totalMarketValue,
    }));
});

const monthlyPerformanceChartOptions = computed(() => {
    const names = monthlyData.value.map(d => d.name);
    const values = monthlyData.value.map(d => d.value / 10000);

    return {
        tooltip: {
            trigger: 'axis',
            axisPointer: { type: 'shadow' }
        },
        grid: {
            left: '3%', right: '4%', bottom: '3%', top: '10%',
            containLabel: true
        },
        xAxis: {
            type: 'category',
            data: names,
            axisLine: {
                lineStyle: { color: isDarkMode.value ? '#424242' : '#e0e0e0' }
            },
            axisLabel: {
                color: isDarkMode.value ? '#bdbdbd' : '#757575',
                rotate: names.length > 6 ? 45 : 0,
            }
        },
        yAxis: {
            type: 'value',
            position: 'left',
            axisLine: { show: false },
            axisTick: { show: false },
            splitLine: {
                lineStyle: {
                    color: isDarkMode.value ? '#424242' : '#e0e0e0',
                    type: 'dashed'
                }
            },
            axisLabel: {
                color: isDarkMode.value ? '#bdbdbd' : '#757575'
            }
        },
        series: [
            {
                type: 'bar',
                data: values,
                itemStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: 'rgba(198, 126, 72, 0.8)' },
                        { offset: 1, color: 'rgba(198, 126, 72, 0.3)' }
                    ])
                },
                barWidth: '50%'
            }
        ]
    };
});

async function reload(force: boolean): Promise<void> {
    loading.value = true;
    try {
        await Promise.all([
            investmentStore.loadOverview({ force }),
            investmentStore.loadHoldings({ force }),
        ]);
        if (force) {
            snackbar.value?.showMessage('Data has been updated');
        }
    } finally {
        loading.value = false;
    }
}

onMounted(() => {
    reload(false);
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

.asset-allocation-chart {
    width: 100%;
    height: 320px;
}

.monthly-performance-chart {
    width: 100%;
    height: 320px;
}

.v-card {
    transition: all 0.3s ease;
}

.v-card:hover {
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}
</style>

<template>
    <div class="page-content">
        <div class="page-header">
            <h1 class="page-title">{{ tt('Investment Overview') }}</h1>
        </div>
        <div class="page-body">
            <v-row class="match-height">
                <v-col cols="12" lg="6" md="12">
                    <v-card :class="{ 'disabled': loading }">
                        <template #title>
                            <div class="d-flex align-center">
                                <span class="text-xl font-weight-bold">{{ tt('Asset Allocation') }}</span>
                                <v-btn density="compact" color="default" variant="text" size="24"
                                       class="ms-2" :icon="true" :loading="loading" @click="reload(true)">
                                    <template #loader>
                                        <v-progress-circular indeterminate size="20"/>
                                    </template>
                                    <v-icon :icon="mdiRefresh" size="24" />
                                    <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                                </v-btn>
                            </div>
                        </template>

                        <v-card-text>
                            <div class="d-flex flex-column align-center">
                                <v-chart autoresize class="asset-allocation-chart" :option="assetAllocationChartOptions" />
                            </div>
                            <div class="mt-4">
                                <div class="d-flex align-center mb-2" v-for="(item, index) in assetAllocationData" :key="index">
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
                                        <span class="text-h5 font-weight-bold" v-if="!loading">{{ portfolioSummary.totalInvestment }}</span>
                                        <v-skeleton-loader class="mt-1" width="100px" type="text" :loading="true" v-else></v-skeleton-loader>
                                    </div>
                                </v-col>

                                <v-col cols="12" md="6">
                                    <div class="d-flex flex-column">
                                        <span class="text-caption text-medium-emphasis">{{ tt('Current Value') }}</span>
                                        <span class="text-h5 font-weight-bold" v-if="!loading">{{ portfolioSummary.currentValue }}</span>
                                        <v-skeleton-loader class="mt-1" width="100px" type="text" :loading="true" v-else></v-skeleton-loader>
                                    </div>
                                </v-col>
                            </v-row>

                            <v-divider class="my-4" />

                            <v-row>
                                <v-col cols="12" md="6">
                                    <div class="d-flex flex-column">
                                        <span class="text-caption text-medium-emphasis">{{ tt('Annualized Return') }}</span>
                                        <div class="d-flex align-center">
                                            <span class="text-h6 font-weight-bold" :class="getReturnColorClass(portfolioSummary.annualizedReturn)" v-if="!loading">{{ portfolioSummary.annualizedReturn }}</span>
                                            <v-select
                                                v-model="selectedReturnType"
                                                :items="returnTypes"
                                                item-title="label"
                                                item-value="value"
                                                density="compact"
                                                variant="outlined"
                                                hide-details
                                                class="return-type-select ml-2"
                                            />
                                        </div>
                                    </div>
                                </v-col>

                                <v-col cols="12" md="6">
                                    <div class="d-flex flex-column">
                                        <span class="text-caption text-medium-emphasis">{{ tt('Cumulative Return') }}</span>
                                        <span class="text-h6 font-weight-bold" :class="getReturnColorClass(portfolioSummary.cumulativeReturn)" v-if="!loading">{{ portfolioSummary.cumulativeReturn }}</span>
                                        <v-skeleton-loader class="mt-1" width="80px" type="text" :loading="true" v-else></v-skeleton-loader>
                                    </div>
                                </v-col>
                            </v-row>

                            <v-divider class="my-4" />

                            <div class="d-flex align-center justify-space-between">
                                <div class="d-flex flex-column">
                                    <span class="text-caption text-medium-emphasis">{{ tt('Total Return') }}</span>
                                    <span class="text-h5 font-weight-bold" :class="getReturnColorClass(portfolioSummary.totalReturn)" v-if="!loading">{{ portfolioSummary.totalReturn }}</span>
                                    <v-skeleton-loader class="mt-1" width="80px" type="text" :loading="true" v-else></v-skeleton-loader>
                                </div>
                                <div class="text-end">
                                    <v-chip :color="getReturnColorClass(portfolioSummary.todayReturn, true)" size="small" variant="tonal">
                                        {{ tt('Today') }}: {{ portfolioSummary.todayReturn }}
                                    </v-chip>
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
                                :return-amount="todayReturn"
                                :datetime="displayDateRange?.today?.displayTime || ''"
                            />
                        </v-col>

                        <v-col cols="6">
                            <investment-return-overview-card
                                :loading="loading" :disabled="loading" :icon="mdiCalendarWeekOutline"
                                :title="tt('This Week')"
                                :return-amount="weekReturn"
                                :datetime="displayDateRange?.thisWeek?.startTime + '-' + displayDateRange?.thisWeek?.endTime"
                            />
                        </v-col>

                        <v-col cols="6">
                            <investment-return-overview-card
                                :loading="loading" :disabled="loading" :icon="mdiCalendarMonthOutline"
                                :title="tt('This Month')"
                                :return-amount="monthReturn"
                                :datetime="displayDateRange?.thisMonth?.startTime + '-' + displayDateRange?.thisMonth?.endTime"
                            />
                        </v-col>

                        <v-col cols="6">
                            <investment-return-overview-card
                                :loading="loading" :disabled="loading" :icon="mdiLayersTripleOutline"
                                :title="tt('This Year')"
                                :return-amount="yearReturn"
                                :datetime="displayDateRange?.thisYear?.displayTime || ''"
                            />
                        </v-col>
                    </v-row>
                </v-col>

                <v-col cols="12" md="6">
                    <v-card :class="{ 'disabled': loading }">
                        <template #title>
                            <div class="d-flex align-center justify-space-between w-100">
                                <span class="text-xl font-weight-bold">{{ tt('Monthly Performance') }}</span>
                                <v-select
                                    v-model="selectedChartType"
                                    :items="chartTypes"
                                    item-title="label"
                                    item-value="value"
                                    density="compact"
                                    variant="outlined"
                                    hide-details
                                    class="chart-type-select"
                                />
                            </div>
                        </template>

                        <v-card-text>
                            <v-chart autoresize class="monthly-performance-chart" :option="monthlyPerformanceChartOptions" />
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

import { ref, computed, useTemplateRef } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import { ThemeType } from '@/core/theme.ts';

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

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const loading = ref<boolean>(true);
const selectedReturnType = ref<string>('cumulative');
const selectedChartType = ref<string>('return');

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

const returnTypes = [
    { label: tt('Cumulative Return'), value: 'cumulative' },
    { label: tt('Time-Weighted Return'), value: 'timeWeighted' },
    { label: tt('Money-Weighted Return'), value: 'moneyWeighted' }
];

const chartTypes = [
    { label: tt('Monthly Return'), value: 'return' },
    { label: tt('Cumulative Return'), value: 'cumulative' },
    { label: tt('Portfolio Value'), value: 'value' },
    { label: tt('Holdings Count'), value: 'holdings' }
];

const displayDateRange = ref({
    today: { displayTime: '2026-03-24' },
    thisWeek: { startTime: '2026-03-17', endTime: '2026-03-23' },
    thisMonth: { startTime: '2026-03-01', endTime: '2026-03-31' },
    thisYear: { displayTime: '2026' }
});

const portfolioSummary = ref({
    totalInvestment: '$100,000',
    currentValue: '$105,000',
    annualizedReturn: '+12.5%',
    cumulativeReturn: '+5.0%',
    totalReturn: '+$5,000',
    todayReturn: '+$120'
});

const todayReturn = computed(() => portfolioSummary.value.todayReturn);
const weekReturn = ref('+$850');
const monthReturn = ref('+$2,300');
const yearReturn = ref('+$5,000');

const assetAllocationData = ref([
    { name: 'Stocks', percent: '60%', value: 63000, color: '#42A5F5' },
    { name: 'ETFs', percent: '20%', value: 21000, color: '#66BB6A' },
    { name: 'Bonds', percent: '15%', value: 15750, color: '#FFA726' },
    { name: 'Cash', percent: '5%', value: 5250, color: '#AB47BC' }
]);

const assetAllocationChartOptions = computed(() => ({
    tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
    },
    legend: {
        show: false
    },
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
            label: {
                show: false,
                position: 'center'
            },
            emphasis: {
                label: {
                    show: true,
                    fontSize: 18,
                    fontWeight: 'bold'
                }
            },
            labelLine: {
                show: false
            },
            data: assetAllocationData.value.map(item => ({
                value: item.value,
                name: item.name,
                itemStyle: {
                    color: item.color
                }
            }))
        }
    ]
}));

const monthlyPerformanceChartOptions = computed(() => {
    let seriesData = [];
    let yAxisConfig = {};
    let seriesConfig = {};

    switch (selectedChartType.value) {
        case 'return':
            seriesData = [2.5, 3.8, 5.0, 6.2, 8.5, 7.1, 9.3, 11.2, 10.1, 12.5, 14.8, 16.2];
            yAxisConfig = {
                type: 'value',
                name: '%',
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
                    formatter: '{value}%',
                    color: isDarkMode.value ? '#bdbdbd' : '#757575'
                }
            };
            seriesConfig = {
                type: 'bar',
                data: seriesData,
                itemStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: 'rgba(66, 165, 245, 0.8)' },
                        { offset: 1, color: 'rgba(66, 165, 245, 0.3)' }
                    ])
                },
                barWidth: '50%'
            };
            break;

        case 'cumulative':
            seriesData = [2.5, 6.4, 11.5, 17.8, 26.5, 33.8, 43.2, 54.5, 64.7, 77.3, 92.2, 108.5];
            yAxisConfig = {
                type: 'value',
                name: '%',
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
                    formatter: '{value}%',
                    color: isDarkMode.value ? '#bdbdbd' : '#757575'
                }
            };
            seriesConfig = {
                type: 'line',
                data: seriesData,
                smooth: true,
                lineStyle: {
                    width: 3,
                    color: '#66BB6A'
                },
                itemStyle: {
                    color: '#66BB6A'
                },
                areaStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: 'rgba(102, 187, 106, 0.4)' },
                        { offset: 1, color: 'rgba(102, 187, 106, 0.05)' }
                    ])
                },
                symbol: 'circle',
                symbolSize: 6
            };
            break;

        case 'value':
            seriesData = [95000, 98000, 100000, 102000, 105000, 103000, 107000, 110000, 108000, 112000, 115000, 118000];
            yAxisConfig = {
                type: 'value',
                name: '$',
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
                    formatter: '${value}',
                    color: isDarkMode.value ? '#bdbdbd' : '#757575'
                }
            };
            seriesConfig = {
                type: 'line',
                data: seriesData,
                smooth: true,
                lineStyle: {
                    width: 3,
                    color: '#42A5F5'
                },
                itemStyle: {
                    color: '#42A5F5'
                },
                areaStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: 'rgba(66, 165, 245, 0.4)' },
                        { offset: 1, color: 'rgba(66, 165, 245, 0.05)' }
                    ])
                },
                symbol: 'circle',
                symbolSize: 6
            };
            break;

        case 'holdings':
            seriesData = [5, 6, 7, 7, 8, 9, 10, 10, 11, 12, 12, 13];
            yAxisConfig = {
                type: 'value',
                name: '',
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
            };
            seriesConfig = {
                type: 'bar',
                data: seriesData,
                itemStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: 'rgba(171, 71, 188, 0.8)' },
                        { offset: 1, color: 'rgba(171, 71, 188, 0.3)' }
                    ])
                },
                barWidth: '50%'
            };
            break;
    }

    return {
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            }
        },
        grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            top: '10%',
            containLabel: true
        },
        xAxis: {
            type: 'category',
            data: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
            axisLine: {
                lineStyle: {
                    color: isDarkMode.value ? '#424242' : '#e0e0e0'
                }
            },
            axisLabel: {
                color: isDarkMode.value ? '#bdbdbd' : '#757575'
            }
        },
        yAxis: yAxisConfig,
        series: [seriesConfig]
    };
});

function getReturnColorClass(returnValue: string, isChip: boolean = false): string {
    if (returnValue.startsWith('+')) {
        return 'text-success';
    } else if (returnValue.startsWith('-')) {
        return 'text-error';
    }
    return '';
}

function reload(force: boolean): void {
    loading.value = true;

    setTimeout(() => {
        loading.value = false;

        if (force) {
            snackbar.value?.showMessage('Data has been updated');
        }
    }, 1000);
}

if (true) {
    reload(false);
}
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

.asset-allocation-chart {
    width: 100%;
    height: 320px;
}

.monthly-performance-chart {
    width: 100%;
    height: 320px;
}

.return-type-select {
    max-width: 160px;
    font-size: 0.875rem;
}

.return-type-select :deep(.v-field) {
    font-size: 0.875rem;
}

.chart-type-select {
    max-width: 160px;
}

.chart-type-select :deep(.v-field) {
    font-size: 0.875rem;
}

.v-card {
    transition: all 0.3s ease;
}

.v-card:hover {
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}
</style>

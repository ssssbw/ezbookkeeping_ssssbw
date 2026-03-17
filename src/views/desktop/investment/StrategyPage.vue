<template>
    <div class="page-content">
        <div class="page-header">
            <h1 class="page-title">{{ tt('Strategy Configuration') }}</h1>
        </div>
        <div class="page-body">
            <v-card class="mb-4">
                <v-card-title>{{ tt('Investment Strategy') }}</v-card-title>
                <v-card-text>
                    <v-form>
                        <v-container>
                            <v-row>
                                <v-col cols="12" md="6">
                                    <v-select
                                        label="Strategy Type"
                                        v-model="strategyType"
                                        :items="strategyTypes"
                                        variant="outlined"
                                        class="mb-4"
                                    />
                                </v-col>
                                <v-col cols="12" md="6">
                                    <v-select
                                        label="Risk Level"
                                        v-model="riskLevel"
                                        :items="riskLevels"
                                        variant="outlined"
                                        class="mb-4"
                                    />
                                </v-col>
                            </v-row>
                            <v-row>
                                <v-col cols="12">
                                    <v-textarea
                                        label="Strategy Description"
                                        v-model="strategyDescription"
                                        variant="outlined"
                                        rows="3"
                                        class="mb-4"
                                    />
                                </v-col>
                            </v-row>
                            <v-row>
                                <v-col cols="12" md="4">
                                    <v-text-field
                                        label="Stock Allocation (%)"
                                        v-model="stockAllocation"
                                        type="number"
                                        variant="outlined"
                                        class="mb-4"
                                    />
                                </v-col>
                                <v-col cols="12" md="4">
                                    <v-text-field
                                        label="Bond Allocation (%)"
                                        v-model="bondAllocation"
                                        type="number"
                                        variant="outlined"
                                        class="mb-4"
                                    />
                                </v-col>
                                <v-col cols="12" md="4">
                                    <v-text-field
                                        label="Cash Allocation (%)"
                                        v-model="cashAllocation"
                                        type="number"
                                        variant="outlined"
                                        class="mb-4"
                                    />
                                </v-col>
                            </v-row>
                        </v-container>
                    </v-form>
                </v-card-text>
                <v-card-actions>
                    <v-btn color="primary">{{ tt('Save Strategy') }}</v-btn>
                    <v-btn color="secondary">{{ tt('Load Default') }}</v-btn>
                </v-card-actions>
            </v-card>
            <v-card class="mb-4">
                <v-card-title>{{ tt('Strategy Backtesting') }}</v-card-title>
                <v-card-text>
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                        <v-select
                            label="Time Period"
                            v-model="backtestPeriod"
                            :items="backtestPeriods"
                            variant="outlined"
                        />
                        <v-select
                            label="Benchmark"
                            v-model="benchmark"
                            :items="benchmarks"
                            variant="outlined"
                        />
                        <v-btn color="primary" class="self-end">{{ tt('Run Backtest') }}</v-btn>
                    </div>
                    <div class="h-80">
                        <!-- 回测结果图表 -->
                        <div class="flex items-center justify-center h-full">
                            <div class="text-gray-500">{{ tt('Backtest results chart will be displayed here') }}</div>
                        </div>
                    </div>
                </v-card-text>
            </v-card>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from '@/locales/helpers.ts';

const { tt } = useI18n();

const strategyType = ref('Balanced');
const strategyTypes = ['Conservative', 'Balanced', 'Aggressive', 'Growth', 'Value'];

const riskLevel = ref('Medium');
const riskLevels = ['Low', 'Medium', 'High', 'Very High'];

const strategyDescription = ref('A balanced investment strategy with equal allocation to stocks and bonds.');

const stockAllocation = ref('50');
const bondAllocation = ref('40');
const cashAllocation = ref('10');

const backtestPeriod = ref('1 Year');
const backtestPeriods = ['3 Months', '6 Months', '1 Year', '3 Years', '5 Years'];

const benchmark = ref('S&P 500');
const benchmarks = ['S&P 500', 'NASDAQ', 'Dow Jones', 'Russell 2000'];
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
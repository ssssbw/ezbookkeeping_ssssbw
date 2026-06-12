<template>
    <v-card-text class="py-3 px-6">
        <div class="d-flex flex-wrap align-center ga-6">
            <div class="summary-item">
                <span class="text-caption text-medium-emphasis">{{ tt('asset.TotalMarketValue') }}</span>
                <span class="text-body-1 font-weight-bold ml-2">{{ formatCurrencyValue(summary.totalMarketValue, 'CNY') }}</span>
            </div>
            <div class="summary-item">
                <span class="text-caption text-medium-emphasis">{{ tt('Total Cost') }}</span>
                <span class="text-body-1 font-weight-bold ml-2">{{ formatCurrencyValue(summary.totalCost, 'CNY') }}</span>
            </div>
            <div class="summary-item">
                <span class="text-caption text-medium-emphasis">{{ tt('Unrealized P&L') }}</span>
                <span class="text-body-1 font-weight-bold ml-2" :class="getReturnColorClass(summary.totalUnrealizedPnl)">
                    {{ formatCurrencyValue(summary.totalUnrealizedPnl, 'CNY') }}
                </span>
            </div>
            <div class="summary-item">
                <span class="text-caption text-medium-emphasis">{{ tt('Return Rate') }}</span>
                <span class="text-body-1 font-weight-bold ml-2" :class="getReturnColorClass(summary.totalReturnRate)">
                    {{ formatReturnRate(summary.totalReturnRate) }}
                </span>
            </div>
        </div>
    </v-card-text>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';
import { formatCurrencyValue, formatReturnRate, getReturnColorClass } from './assetUtils.ts';

const { tt } = useI18n();
const investmentStore = useInvestmentStore();

const summary = investmentStore.holdingsSummary;
</script>

<style scoped>
.summary-item {
    display: flex;
    align-items: center;
    white-space: nowrap;
}
</style>

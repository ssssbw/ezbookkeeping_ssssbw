<template>
    <v-card-text class="py-4 px-6">
        <v-row dense>
            <v-col cols="12" sm="6" md="3">
                <div class="d-flex flex-column">
                    <span class="text-caption text-medium-emphasis">{{ tt('asset.TotalMarketValue') }}</span>
                    <span class="text-h6 font-weight-bold mt-1">{{ formatCurrencyValue(summary.totalMarketValue, 'CNY') }}</span>
                </div>
            </v-col>
            <v-col cols="12" sm="6" md="3">
                <div class="d-flex flex-column">
                    <span class="text-caption text-medium-emphasis">{{ tt('Total Cost') }}</span>
                    <span class="text-h6 font-weight-bold mt-1">{{ formatCurrencyValue(summary.totalCost, 'CNY') }}</span>
                </div>
            </v-col>
            <v-col cols="12" sm="6" md="3">
                <div class="d-flex flex-column">
                    <span class="text-caption text-medium-emphasis">{{ tt('Unrealized P&L') }}</span>
                    <span class="text-h6 font-weight-bold mt-1" :class="getReturnColorClass(summary.totalUnrealizedPnl)">
                        {{ formatCurrencyValue(summary.totalUnrealizedPnl, 'CNY') }}
                    </span>
                </div>
            </v-col>
            <v-col cols="12" sm="6" md="3">
                <div class="d-flex flex-column">
                    <span class="text-caption text-medium-emphasis">{{ tt('Return Rate') }}</span>
                    <span class="text-h6 font-weight-bold mt-1" :class="getReturnColorClass(summary.totalReturnRate)">
                        {{ formatReturnRate(summary.totalReturnRate) }}
                    </span>
                </div>
            </v-col>
        </v-row>
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
/* All styling uses Vuetify utility classes, consistent with OverviewPage */
</style>

<template>
    <v-dialog :model-value="modelValue" max-width="600" @update:model-value="$emit('update:modelValue', $event)">
        <v-card v-if="item">
            <v-card-title class="d-flex align-center">
                <span class="text-h6">{{ item.assetName }}</span>
                <v-spacer />
                <v-chip size="small" variant="tonal" :color="getCategoryColor(item.category)">
                    {{ formatCategory(item.category, tt) }}
                </v-chip>
                <v-btn variant="text" :icon="mdiClose" density="compact" @click="$emit('update:modelValue', false)" />
            </v-card-title>
            <v-card-text>
                <v-row dense>
                    <v-col cols="6">
                        <div class="text-caption text-medium-emphasis">{{ tt('Code') }}</div>
                        <div class="text-body-1">{{ item.assetCode }}</div>
                    </v-col>
                    <v-col cols="6">
                        <div class="text-caption text-medium-emphasis">{{ tt('Market') }}</div>
                        <div class="text-body-1">{{ formatMarket(item.market, tt) }}</div>
                    </v-col>
                    <v-col cols="6">
                        <div class="text-caption text-medium-emphasis">{{ tt('Category') }}</div>
                        <div class="text-body-1">{{ formatCategory(item.category, tt) }}</div>
                    </v-col>
                    <v-col cols="6">
                        <div class="text-caption text-medium-emphasis">{{ tt('Currency') }}</div>
                        <div class="text-body-1">{{ item.currency }}</div>
                    </v-col>
                </v-row>

                <template v-if="aggregated">
                    <v-divider class="my-3" />
                    <v-row dense>
                        <v-col cols="6">
                            <div class="text-caption text-medium-emphasis">{{ tt('Current Price') }}</div>
                            <div class="text-body-1 font-weight-medium">
                                {{ formatPriceWithDate(aggregated.holdings[0]?.currentPrice, aggregated.holdings[0]?.currentPriceDate) }}
                            </div>
                        </v-col>
                        <v-col cols="6">
                            <div class="text-caption text-medium-emphasis">{{ tt('Total Quantity') }}</div>
                            <div class="text-body-1">{{ formatQuantity(aggregated.totalQuantity) }}</div>
                        </v-col>
                        <v-col cols="6">
                            <div class="text-caption text-medium-emphasis">{{ tt('Total Cost') }}</div>
                            <div class="text-body-1">{{ formatCurrencyValue(aggregated.totalCost, item.currency) }}</div>
                        </v-col>
                        <v-col cols="6">
                            <div class="text-caption text-medium-emphasis">{{ tt('Market Value') }}</div>
                            <div class="text-body-1 font-weight-medium">{{ formatCurrencyValue(aggregated.totalMarketValue, item.currency) }}</div>
                        </v-col>
                        <v-col cols="6">
                            <div class="text-caption text-medium-emphasis">{{ tt('Unrealized P&L') }}</div>
                            <div class="text-body-1" :class="getReturnColorClass(aggregated.unrealizedPnl)">
                                {{ formatCurrencyValue(aggregated.unrealizedPnl, item.currency) }}
                            </div>
                        </v-col>
                        <v-col cols="6">
                            <div class="text-caption text-medium-emphasis">{{ tt('Return Rate') }}</div>
                            <div class="text-body-1 font-weight-bold" :class="getReturnColorClass(aggregated.weightedReturnRate)">
                                {{ formatReturnRate(aggregated.weightedReturnRate) }}
                            </div>
                        </v-col>
                    </v-row>
                </template>

                <template v-else-if="item.currentPrice !== undefined && item.currentPrice !== null">
                    <v-divider class="my-3" />
                    <v-row dense>
                        <v-col cols="12">
                            <div class="text-caption text-medium-emphasis">{{ tt('Current Price') }}</div>
                            <div class="text-body-1">{{ formatPriceWithDate(item.currentPrice, item.currentPriceDate) }}</div>
                        </v-col>
                    </v-row>
                </template>
            </v-card-text>
            <v-card-actions>
                <v-spacer />
                <v-btn variant="text" @click="$emit('update:modelValue', false)">{{ tt('Close') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';

import { mdiClose } from '@mdi/js';

import type { DisplayAsset } from './types.ts';

import {
    formatMarket, formatCategory, formatPriceWithDate, formatQuantity,
    formatCurrencyValue, formatReturnRate, getReturnColorClass, getCategoryColor
} from './assetUtils.ts';

const props = defineProps<{
    modelValue: boolean;
    item: DisplayAsset | null;
}>();

defineEmits<{
    (e: 'update:modelValue', value: boolean): void;
}>();

const { tt } = useI18n();
const investmentStore = useInvestmentStore();

const aggregated = computed(() => {
    if (!props.item) return null;
    return investmentStore.aggregatedHoldings.find(a => a.assetId === props.item!.assetId) || null;
});
</script>

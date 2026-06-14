<template>
    <v-card-text class="pa-0">
        <v-data-table
            :headers="headers"
            :items="aggregatedHoldings"
            :loading="loading"
            :hover="true"
            item-value="assetId"
            show-expand
            class="holdings-table"
            v-model:items-per-page="perPage"
            v-model:page="page"
            @click:row="onRowClick"
        >
            <template #item.assetCode="{ item }">
                <span class="text-body-2 font-weight-medium">{{ item.assetCode }}</span>
            </template>
            <template #item.assetName="{ item }">
                <span class="text-body-2">{{ item.assetName }}</span>
            </template>
            <template #item.market="{ item }">
                <span class="text-body-2">{{ formatMarket(item.market, tt) }}</span>
            </template>
            <template #item.totalQuantity="{ item }">
                <span class="text-body-2">{{ formatQuantity(item.totalQuantity) }}</span>
            </template>
            <template #item.totalMarketValue="{ item }">
                <span class="text-body-2 font-weight-medium">{{ formatCurrencyValue(item.totalMarketValue, item.currency) }}</span>
            </template>
            <template #item.totalCost="{ item }">
                <span class="text-body-2">{{ formatCurrencyValue(item.totalCost, item.currency) }}</span>
            </template>
            <template #item.unrealizedPnl="{ item }">
                <span class="text-body-2 font-weight-bold" :class="getReturnColorClass(item.unrealizedPnl)">
                    {{ formatCurrencyValue(item.unrealizedPnl, item.currency) }}
                </span>
            </template>
            <template #item.weightedReturnRate="{ item }">
                <span class="text-body-2 font-weight-bold" :class="getReturnColorClass(item.weightedReturnRate)">
                    {{ formatReturnRate(item.weightedReturnRate) }}
                </span>
            </template>

            <!-- Expanded row: per-account breakdown (single td colspan, grid layout, robust to column changes) -->
            <template #expanded-row="{ item, columns }">
                <tr>
                    <td :colspan="columns.length" class="pa-0">
                        <div class="expanded-detail">
                            <div v-for="h in item.holdings" :key="h.accountId" class="expanded-detail-row">
                                <span class="text-body-2 text-medium-emphasis account-label">
                                    {{ h.accountName || h.accountId }}
                                </span>
                                <span class="text-body-2 detail-value">{{ formatQuantity(h.quantity) }}</span>
                                <span class="text-body-2 detail-value">{{ formatCurrencyValue(h.marketValue, h.currency) }}</span>
                                <span class="text-body-2 detail-value">{{ formatCurrencyValue(h.totalCost, h.currency) }}</span>
                                <span class="text-body-2 detail-value font-weight-medium" :class="getReturnColorClass(h.unrealizedPnl)">
                                    {{ formatCurrencyValue(h.unrealizedPnl, h.currency) }}
                                </span>
                                <span class="text-body-2 detail-value" :class="getReturnColorClass(h.returnRate)">
                                    {{ formatReturnRate(h.returnRate) }}
                                </span>
                            </div>
                        </div>
                    </td>
                </tr>
            </template>

            <template #loading>
                <v-skeleton-loader type="table-row@10" :loading="true" />
            </template>
            <template #no-data>
                <div class="text-center py-6 text-medium-emphasis">
                    {{ tt('No holdings data') }}
                </div>
            </template>
            <template #bottom>
                <div class="mt-2 mb-4">
                    <pagination-buttons :totalPageCount="totalPageCount" v-model="page" />
                </div>
            </template>
        </v-data-table>
    </v-card-text>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import PaginationButtons from '@/components/desktop/PaginationButtons.vue';

import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';

import {
    formatMarket, formatQuantity, formatCurrencyValue,
    formatReturnRate, getReturnColorClass
} from './assetUtils.ts';

const { tt } = useI18n();
const investmentStore = useInvestmentStore();

defineProps<{
    loading: boolean;
}>();

const perPage = ref<number>(10);
const page = ref<number>(1);

const aggregatedHoldings = computed(() => investmentStore.aggregatedHoldings);

const totalPageCount = computed<number>(() => {
    return Math.max(1, Math.ceil(aggregatedHoldings.value.length / (perPage.value > 0 ? perPage.value : 10)));
});

const headers = computed(() => [
    { key: 'assetCode', title: tt('Code'), sortable: false },
    { key: 'assetName', title: tt('asset.AssetName'), sortable: false },
    { key: 'market', title: tt('Market'), sortable: false },
    { key: 'totalQuantity', title: tt('Total Quantity'), sortable: false, align: 'end' as const },
    { key: 'totalMarketValue', title: tt('Market Value'), sortable: false, align: 'end' as const },
    { key: 'totalCost', title: tt('Total Cost'), sortable: false, align: 'end' as const },
    { key: 'unrealizedPnl', title: tt('Unrealized P&L'), sortable: false, align: 'end' as const },
    { key: 'weightedReturnRate', title: tt('Return Rate'), sortable: false, align: 'end' as const },
]);

function onRowClick(_event: Event, { internalItem, toggleExpand }: { internalItem: unknown; toggleExpand: (item: unknown) => void }): void {
    toggleExpand(internalItem);
}
</script>

<style scoped>
.holdings-table :deep(.v-data-table__td) {
    padding-top: 8px;
    padding-bottom: 8px;
}

.expanded-detail {
    background: rgba(var(--v-theme-on-surface), 0.04);
    border-inline-start: 3px solid rgb(var(--v-theme-primary));
}

.expanded-detail-row {
    display: grid;
    grid-template-columns: 1fr repeat(5, minmax(80px, 1fr));
    align-items: center;
    column-gap: 16px;
    padding: 6px 16px 6px 20px;
}

.account-label {
    display: block;
    padding-left: 64px;
}

.detail-value {
    font-variant-numeric: tabular-nums;
    text-align: end;
}
</style>

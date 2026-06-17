<template>
    <v-card-text class="pa-0 holdings-tab">
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
                <span class="text-body-2 asset-code">{{ item.assetCode }}</span>
            </template>
            <template #item.assetName="{ item }">
                <span class="text-body-2">{{ item.assetName }}</span>
            </template>
            <template #item.market="{ item }">
                <span class="text-body-2">{{ formatMarket(item.market, tt) }}</span>
            </template>
            <template #item.currentPrice="{ item }">
                <span class="text-body-2">{{ formatPrice(item.currentPrice) }}</span>
            </template>
            <template #item.totalQuantity="{ item }">
                <span class="text-body-2">{{ formatQuantity(item.totalQuantity) }}</span>
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
            <template #item.date>
                <span class="text-body-2">{{ todayDate }}</span>
            </template>

            <!-- Expanded row: per-account breakdown -->
            <template #expanded-row="{ item, columns }">
                <tr>
                    <td :colspan="columns.length" class="pa-0">
                        <div class="expanded-detail">
                            <div v-for="h in item.holdings" :key="h.accountId" class="expanded-detail-row">
                                <span class="text-body-2 text-medium-emphasis account-label">
                                    {{ h.accountName || h.accountId }}
                                </span>
                                <span class="text-body-2">{{ formatQuantity(h.quantity) }}</span>
                                <span class="text-body-2">{{ formatCurrencyValue(h.marketValue, h.currency) }}</span>
                                <span class="text-body-2">{{ formatCurrencyValue(h.totalCost, h.currency) }}</span>
                                <span class="text-body-2" :class="getReturnColorClass(h.unrealizedPnl)">
                                    {{ formatCurrencyValue(h.unrealizedPnl, h.currency) }}
                                </span>
                                <span class="text-body-2" :class="getReturnColorClass(h.returnRate)">
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
    formatMarket, formatPrice, formatQuantity, formatCurrencyValue,
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

const todayDate = computed(() => {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
});

const totalPageCount = computed<number>(() => {
    return Math.max(1, Math.ceil(aggregatedHoldings.value.length / (perPage.value > 0 ? perPage.value : 10)));
});

const headers = computed(() => [
    { key: 'assetCode', title: tt('Code'), sortable: false, width: '100' },
    { key: 'assetName', title: tt('asset.AssetName'), sortable: false },
    { key: 'market', title: tt('Market'), sortable: false, width: '80' },
    { key: 'currentPrice', title: tt('Current Price'), sortable: false, align: 'end' as const },
    { key: 'totalQuantity', title: tt('Total Quantity'), sortable: false, align: 'end' as const },
    { key: 'totalMarketValue', title: tt('Market Value'), sortable: false, align: 'end' as const },
    { key: 'totalCost', title: tt('Total Cost'), sortable: false, align: 'end' as const },
    { key: 'unrealizedPnl', title: tt('Unrealized P&L'), sortable: false, align: 'end' as const },
    { key: 'weightedReturnRate', title: tt('Return Rate'), sortable: false, align: 'end' as const },
    { key: 'date', title: tt('Date'), sortable: false, align: 'center' as const, width: '110' },
]);

function onRowClick(_event: Event, { internalItem, toggleExpand }: { internalItem: unknown; toggleExpand: (item: unknown) => void }): void {
    toggleExpand(internalItem);
}
</script>

<style scoped>
.holdings-tab {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
}

.holdings-table {
    flex: 1;
}

.holdings-table :deep(.v-data-table) {
    height: 100%;
}

.holdings-table :deep(.v-data-table__td) {
    padding-top: 8px;
    padding-bottom: 8px;
}

.asset-code {
    color: rgb(var(--v-theme-primary));
    font-weight: 500;
}

.expanded-detail {
    background: rgba(var(--v-theme-on-surface), 0.03);
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
    padding-left: 180px;
}
</style>

<template>
    <v-card-text class="pa-0 holdings-tab">
        <v-data-table
            :headers="headers"
            :items="aggregatedHoldings"
            :loading="loading"
            :hover="true"
            item-value="assetId"
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
                <span class="text-body-2">{{ formatPriceWithDate(item.currentPrice, item.currentPriceDate) }}</span>
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
import { useRouter } from 'vue-router';

import PaginationButtons from '@/components/desktop/PaginationButtons.vue';

import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';

import {
    formatMarket, formatPriceWithDate, formatQuantity, formatCurrencyValue,
    formatReturnRate, getReturnColorClass
} from './assetUtils.ts';

const { tt } = useI18n();
const router = useRouter();
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
    { key: 'assetCode', title: tt('asset.AssetCode'), sortable: false, width: '100' },
    { key: 'assetName', title: tt('asset.AssetName'), sortable: false },
    { key: 'market', title: tt('Market'), sortable: false, width: '80' },
    { key: 'currentPrice', title: tt('Current Price'), sortable: false },
    { key: 'totalQuantity', title: tt('Total Quantity'), sortable: false },
    { key: 'totalMarketValue', title: tt('Total Value'), sortable: false },
    { key: 'totalCost', title: tt('Total Cost'), sortable: false },
    { key: 'unrealizedPnl', title: tt('Unrealized P&L'), sortable: false },
    { key: 'weightedReturnRate', title: tt('Return Rate'), sortable: false },
]);

function onRowClick(_event: Event, { item }: { item: { assetId: string } }): void {
    router.push(`/investment/assets/${item.assetId}`);
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

</style>

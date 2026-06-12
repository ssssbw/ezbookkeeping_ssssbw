<template>
    <v-card-text class="pa-0">
        <v-data-table
            :headers="headers"
            :items="watchlist"
            :loading="loading"
            :hover="true"
            item-value="assetId"
            class="watchlist-table"
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
            <template #item.industry="{ item }">
                <span v-if="item.industry" class="text-body-2">{{ formatIndustry(item.industry, tt) }}</span>
                <span v-else class="text-medium-emphasis">--</span>
            </template>
            <template #item.currentPrice="{ item }">
                <span v-if="item.currentPrice !== undefined && item.currentPrice !== null" class="text-body-2">
                    {{ formatPrice(item.currentPrice) }}
                </span>
                <span v-else class="text-medium-emphasis">--</span>
            </template>
            <template #item.actions="{ item }">
                <div class="d-flex ga-1">
                    <v-btn size="x-small" variant="tonal" color="error" @click.stop="removeFromWatchlist(item)">
                        {{ tt('Remove') }}
                    </v-btn>
                </div>
            </template>
            <template #loading>
                <v-skeleton-loader type="table-row@10" :loading="true" />
            </template>
            <template #no-data>
                <div class="text-center py-6 text-medium-emphasis">
                    {{ tt('No watchlist data') }}
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

import type { DisplayAsset } from './types.ts';

import { formatMarket, formatIndustry, formatPrice } from './assetUtils.ts';

import logger from '@/lib/logger.ts';

const props = defineProps<{
    loading: boolean;
    watchlist: DisplayAsset[];
}>();

const emit = defineEmits<{
    (e: 'select', item: DisplayAsset): void;
    (e: 'removed'): void;
}>();

const { tt } = useI18n();
const investmentStore = useInvestmentStore();

const perPage = ref<number>(10);
const page = ref<number>(1);

const totalPageCount = computed<number>(() => {
    return Math.max(1, Math.ceil(props.watchlist.length / (perPage.value > 0 ? perPage.value : 10)));
});

const headers = computed(() => [
    { key: 'assetCode', title: tt('Code'), sortable: false },
    { key: 'assetName', title: tt('asset.AssetName'), sortable: false },
    { key: 'market', title: tt('Market'), sortable: false },
    { key: 'industry', title: tt('Industry'), sortable: false },
    { key: 'currentPrice', title: tt('Current Price'), sortable: false, align: 'end' as const },
    { key: 'actions', title: tt('Action'), sortable: false, align: 'center' as const, width: '100' },
]);

function onRowClick(_event: Event, row: { item: DisplayAsset }): void {
    emit('select', row.item);
}

async function removeFromWatchlist(item: DisplayAsset): Promise<void> {
    try {
        await investmentStore.removeUserAsset({ assetId: item.assetId });
        emit('removed');
    } catch (error) {
        logger.error('Failed to remove asset from watchlist', error);
    }
}
</script>

<style scoped>
.watchlist-table :deep(.v-data-table__td) {
    padding-top: 8px;
    padding-bottom: 8px;
}
</style>

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
                <span class="text-body-2 asset-code">{{ item.assetCode }}</span>
            </template>
            <template #item.assetName="{ item }">
                <span class="text-body-2">{{ item.assetName }}</span>
            </template>
            <template #item.market="{ item }">
                <span class="text-body-2">{{ formatMarket(item.market, tt) }}</span>
            </template>
            <template #item.currentPrice="{ item }">
                <span v-if="item.currentPrice !== undefined && item.currentPrice !== null" class="text-body-2">
                    {{ formatPriceWithDate(item.currentPrice, item.currentPriceDate) }}
                </span>
                <span v-else class="text-medium-emphasis">--</span>
            </template>
            <template #item.changePercent="{ item }">
                <span v-if="item.currentPrice && item.addedPrice" class="text-body-2" :class="getReturnColorClass(calcChangePercent(item))">
                    {{ formatReturnRate(calcChangePercent(item)) }}
                </span>
                <span v-else class="text-medium-emphasis">--</span>
            </template>
            <template #item.date>
                <span class="text-body-2">{{ todayDate }}</span>
            </template>
            <template #item.actions="{ item }">
                <div class="d-flex ga-1 justify-end">
                    <v-btn size="small" variant="text" color="error" density="comfortable" @click.stop="removeFromWatchlist(item)">
                        <v-icon :icon="mdiStarRemove" />
                        <v-tooltip activator="parent" location="top">{{ tt('Remove from Watchlist') }}</v-tooltip>
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
import { useRouter } from 'vue-router';

import PaginationButtons from '@/components/desktop/PaginationButtons.vue';

import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';

import type { DisplayAsset } from './types.ts';

import { formatMarket, formatPriceWithDate, formatReturnRate, getReturnColorClass } from './assetUtils.ts';

import { mdiStarRemove } from '@mdi/js';

import logger from '@/lib/logger.ts';

const props = defineProps<{
    loading: boolean;
    watchlist: DisplayAsset[];
}>();

const emit = defineEmits<{
    (e: 'removed'): void;
}>();

const { tt } = useI18n();
const router = useRouter();
const investmentStore = useInvestmentStore();

const perPage = ref<number>(10);
const page = ref<number>(1);

const todayDate = computed(() => {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
});

const totalPageCount = computed<number>(() => {
    return Math.max(1, Math.ceil(props.watchlist.length / (perPage.value > 0 ? perPage.value : 10)));
});

const headers = computed(() => [
    { key: 'assetCode', title: tt('Code'), sortable: false, width: '100' },
    { key: 'assetName', title: tt('asset.AssetName'), sortable: false },
    { key: 'market', title: tt('Market'), sortable: false, width: '80' },
    { key: 'currentPrice', title: tt('Current Price'), sortable: false },
    { key: 'changePercent', title: tt('Change Percent'), sortable: false },
    { key: 'date', title: tt('Date'), sortable: false, width: '110' },
    { key: 'actions', title: '', sortable: false, width: '60' },
]);

function calcChangePercent(item: DisplayAsset): number {
    if (!item.currentPrice || !item.addedPrice || item.addedPrice === 0) return 0;
    return Math.round(((item.currentPrice - item.addedPrice) / item.addedPrice) * 10000);
}

function onRowClick(_event: Event, row: { item: DisplayAsset }): void {
    router.push(`/investment/assets/${row.item.assetId}`);
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

.asset-code {
    color: rgb(var(--v-theme-primary));
    font-weight: 500;
}
</style>

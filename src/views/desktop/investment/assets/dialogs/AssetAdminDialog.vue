<template>
    <v-dialog :model-value="modelValue" max-width="1000" scrollable @update:model-value="$emit('update:modelValue', $event)">
        <v-card style="height: 70vh;">
            <v-card-title class="d-flex align-center">
                <span class="text-h6">{{ tt('asset.Global Asset Management') }}</span>
                <v-spacer />
                <v-btn variant="text" :icon="mdiClose" density="compact" @click="$emit('update:modelValue', false)" />
            </v-card-title>
            <v-divider />
            <v-card-text class="pa-0">
                <div class="pa-4 pb-2">
                    <div class="d-flex ga-3 align-center">
                        <v-text-field
                            v-model="searchKeyword"
                            :prepend-inner-icon="mdiMagnify"
                            density="compact"
                            :placeholder="tt('asset.SearchAssetsByNameOrCode')"
                            hide-details
                            clearable
                            variant="outlined"
                            style="max-width: 360px"
                            @update:model-value="onSearchChange"
                        />
                        <v-spacer />
                        <v-btn color="primary" variant="tonal" @click="onAddClick">
                            {{ tt('Add') }}
                        </v-btn>
                        <v-btn
                            variant="tonal"
                            color="secondary"
                            :loading="syncing"
                            :disabled="syncing"
                            :prepend-icon="mdiDatabaseSync"
                            @click="onSync"
                        >
                            {{ tt('asset.Sync Assets') }}
                        </v-btn>
                    </div>
                </div>

                <v-data-table
                    fixed-header
                    fixed-footer
                    :headers="headers"
                    :items="assets"
                    :loading="loading"
                    :hover="true"
                    item-value="id"
                    density="compact"
                    class="admin-table"
                    v-model:items-per-page="perPage"
                    v-model:page="page"
                >
                    <template #item.code="{ item }">
                        <span class="text-body-2 font-weight-medium">{{ item.code }}</span>
                    </template>
                    <template #item.name="{ item }">
                        <span class="text-body-2">{{ item.name }}</span>
                    </template>
                    <template #item.market="{ item }">
                        <span class="text-body-2">{{ formatMarket(item.market, tt) }}</span>
                    </template>
                    <template #item.category="{ item }">
                        <span class="text-body-2">{{ formatCategory(item.category, tt) }}</span>
                    </template>
                    <template #item.industry="{ item }">
                        <span v-if="item.industry" class="text-body-2">{{ formatIndustry(item.industry, tt) }}</span>
                        <span v-else class="text-medium-emphasis">--</span>
                    </template>
                    <template #item.subCategory="{ item }">
                        <span v-if="item.subCategory" class="text-body-2">{{ item.subCategory }}</span>
                        <span v-else class="text-medium-emphasis">--</span>
                    </template>
                    <template #item.actions="{ item }">
                        <div class="d-flex ga-1">
                            <v-btn size="x-small" variant="text" :icon="mdiPencil" @click.stop="onEditClick(item)" />
                            <v-btn size="x-small" variant="text" :icon="mdiDelete" color="error" @click.stop="onDeleteClick(item)" />
                        </div>
                    </template>
                    <template #loading>
                        <v-skeleton-loader type="table-row@10" :loading="true" />
                    </template>
                    <template #no-data>
                        <div class="text-center py-6 text-medium-emphasis">{{ tt('asset.NoAssetsFound') }}</div>
                    </template>
                    <template #bottom>
                        <div class="mt-2 mb-4">
                            <pagination-buttons :totalPageCount="totalPageCount" v-model="page" />
                        </div>
                    </template>
                </v-data-table>
            </v-card-text>
            <v-divider />
            <v-card-actions>
                <v-spacer />
                <v-btn variant="text" @click="$emit('update:modelValue', false)">{{ tt('Close') }}</v-btn>
            </v-card-actions>
        </v-card>

        <AssetFormDialog
            v-model="formDialogVisible"
            :mode="formMode"
            :asset="formAsset"
            @saved="loadAssets"
        />

        <confirm-dialog ref="confirmDialog" />
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, useTemplateRef } from 'vue';

import PaginationButtons from '@/components/desktop/PaginationButtons.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import AssetFormDialog from './AssetFormDialog.vue';

import { useI18n } from '@/locales/helpers.ts';

import type { AssetInfoResponse } from '@/models/investment.ts';

import { mdiMagnify, mdiDatabaseSync, mdiClose, mdiPencil, mdiDelete } from '@mdi/js';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

import { formatMarket, formatCategory, formatIndustry } from '../assetUtils.ts';

const props = defineProps<{
    modelValue: boolean;
}>();

defineEmits<{
    (e: 'update:modelValue', value: boolean): void;
}>();

const { tt } = useI18n();

const loading = ref<boolean>(false);
const syncing = ref<boolean>(false);
const searchKeyword = ref<string>('');
const assets = ref<AssetInfoResponse[]>([]);
const page = ref<number>(1);
const perPage = ref<number>(15);

const formDialogVisible = ref<boolean>(false);
const formMode = ref<'create' | 'edit'>('create');
const formAsset = ref<AssetInfoResponse | null>(null);

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');

const headers = [
    { key: 'code', title: tt('asset.AssetCode'), sortable: false },
    { key: 'name', title: tt('asset.AssetName'), sortable: false },
    { key: 'market', title: tt('Market'), sortable: false },
    { key: 'category', title: tt('Category'), sortable: false },
    { key: 'industry', title: tt('Industry'), sortable: false },
    { key: 'subCategory', title: tt('asset.SubCategory'), sortable: false },
    { key: 'currency', title: tt('Currency'), sortable: false },
    { key: 'actions', title: '', sortable: false, width: '100' }
];

const totalPageCount = computed(() => {
    return Math.max(1, Math.ceil(assets.value.length / perPage.value));
});

watch(() => props.modelValue, (val) => {
    if (val) {
        searchKeyword.value = '';
        loadAssets();
    }
});

let searchTimer: ReturnType<typeof setTimeout> | null = null;

function onSearchChange(): void {
    if (searchTimer) {
        clearTimeout(searchTimer);
    }
    searchTimer = setTimeout(() => {
        loadAssets();
    }, 300);
}

async function loadAssets(): Promise<void> {
    loading.value = true;
    try {
        const resp = await services.listGlobalAssets({
            keyword: searchKeyword.value || undefined
        });
        assets.value = resp.data?.result?.assets || [];
    } catch (e) {
        logger.error('Failed to load admin assets', e);
        assets.value = [];
    } finally {
        loading.value = false;
    }
}

async function onSync(): Promise<void> {
    syncing.value = true;
    try {
        await services.syncGlobalAssets();
        await loadAssets();
    } catch (e) {
        logger.error('Failed to sync assets', e);
    } finally {
        syncing.value = false;
    }
}

function onAddClick(): void {
    formMode.value = 'create';
    formAsset.value = null;
    formDialogVisible.value = true;
}

function onEditClick(asset: AssetInfoResponse): void {
    formMode.value = 'edit';
    formAsset.value = asset;
    formDialogVisible.value = true;
}

async function onDeleteClick(asset: AssetInfoResponse): Promise<void> {
    try {
        await confirmDialog.value?.open(
            tt('Confirm'),
            tt('asset.ConfirmDeleteAsset') + ' "' + asset.name + '" ?'
        );
        await services.deleteGlobalAsset({ id: asset.id });
        await loadAssets();
    } catch {
        // user cancelled or delete failed (error shown by snackbar if needed)
    }
}
</script>

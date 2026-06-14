<template>
    <div class="search-section">
        <div class="search-wrapper">
            <v-text-field
                v-model="searchQuery"
                :placeholder="tt('asset.SearchAssetsByNameOrCode')"
                :prepend-inner-icon="mdiMagnify"
                clearable
                hide-details
                density="compact"
                variant="outlined"
                class="search-input"
                @focus="onSearchFocus"
                @blur="onSearchBlur"
                @update:model-value="onSearchInput"
            />
            <v-card
                v-if="showSearchResults && searchResults.length > 0"
                class="search-results-overlay"
                elevation="8"
            >
                <v-list density="compact" lines="two">
                    <v-list-item
                        v-for="item in searchResults"
                        :key="item.id"
                        @mousedown.prevent="onSearchResultClick(item)"
                    >
                        <template #prepend>
                            <v-avatar size="32" :color="getCategoryColor(item.category)" variant="tonal">
                                <span class="text-caption font-weight-bold">{{ item.code.substring(0, 2) }}</span>
                            </v-avatar>
                        </template>
                        <v-list-item-title class="d-flex align-center">
                            <span class="text-body-2 font-weight-medium">{{ item.code }}</span>
                            <span class="text-body-2 text-medium-emphasis ml-2">{{ item.name }}</span>
                        </v-list-item-title>
                        <v-list-item-subtitle>
                            <span class="text-caption text-medium-emphasis">
                                {{ formatCategory(item.category, tt) }} · {{ formatMarket(item.market, tt) }}
                            </span>
                        </v-list-item-subtitle>
                        <template #append>
                            <v-btn
                                v-if="!isInWatchlist(item.id)"
                                size="small"
                                variant="tonal"
                                color="primary"
                                :icon="mdiStarPlus"
                                @mousedown.prevent.stop="addToWatchlist(item)"
                            />
                        </template>
                    </v-list-item>
                </v-list>
            </v-card>
        </div>
        <v-btn
            v-if="showManageButton"
            variant="tonal"
            color="primary"
            :prepend-icon="mdiCogOutline"
            class="manage-btn ml-3"
            @click="$emit('manage')"
        >
            {{ tt('asset.Global Asset Management') }}
        </v-btn>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';

import { type AssetInfoResponse } from '@/models/investment.ts';

import { mdiMagnify, mdiStarPlus, mdiCogOutline } from '@mdi/js';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

import { formatMarket, formatCategory, getCategoryColor } from './assetUtils.ts';

const props = defineProps<{
    showManageButton: boolean;
    watchlistAssetIds: Set<string>;
}>();

const emit = defineEmits<{
    (e: 'manage'): void;
    (e: 'select', item: AssetInfoResponse): void;
    (e: 'watchlist-updated'): void;
}>();

const { tt } = useI18n();
const investmentStore = useInvestmentStore();

const searchQuery = ref<string>('');
const searchResults = ref<AssetInfoResponse[]>([]);
const showSearchResults = ref<boolean>(false);
let searchTimer: ReturnType<typeof setTimeout> | null = null;

function isInWatchlist(assetId: string): boolean {
    return props.watchlistAssetIds.has(assetId);
}

function onSearchFocus(): void {
    if (searchQuery.value.trim() && searchResults.value.length > 0) {
        showSearchResults.value = true;
    }
}

function onSearchBlur(): void {
    setTimeout(() => {
        showSearchResults.value = false;
    }, 200);
}

function onSearchInput(value: string): void {
    if (!value || !value.trim()) {
        searchResults.value = [];
        showSearchResults.value = false;
        return;
    }

    if (searchTimer) {
        clearTimeout(searchTimer);
    }

    searchTimer = setTimeout(() => {
        performSearch(value.trim());
    }, 300);
}

async function performSearch(keyword: string): Promise<void> {
    try {
        const response = await services.searchAssets({ keyword, limit: 8 });
        const data = response.data;

        if (!data || !data.success || !data.result) {
            searchResults.value = [];
            showSearchResults.value = false;
            return;
        }

        searchResults.value = data.result;
        showSearchResults.value = data.result.length > 0;
    } catch (error) {
        logger.error('Failed to search assets', error);
        searchResults.value = [];
        showSearchResults.value = false;
    }
}

function onSearchResultClick(item: AssetInfoResponse): void {
    showSearchResults.value = false;
    searchQuery.value = '';
    emit('select', item);
}

async function addToWatchlist(item: AssetInfoResponse): Promise<void> {
    try {
        await investmentStore.addUserAsset({ assetId: item.id });
        emit('watchlist-updated');
    } catch (error) {
        logger.error('Failed to add to watchlist', error);
    }
}
</script>

<style scoped>
.search-section {
    display: flex;
    align-items: center;
}

.search-wrapper {
    flex: 1;
    position: relative;
}

.search-input {
    max-width: 100%;
}

.manage-btn {
    flex-shrink: 0;
    height: 40px;
}

.search-results-overlay {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    z-index: 100;
    margin-top: 4px;
    min-height: 200px;
    max-height: 400px;
    overflow-y: auto;
}
</style>

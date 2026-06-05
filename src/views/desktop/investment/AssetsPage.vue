<template>
    <div class="page-content">
        <div class="page-body">
            <!-- Top Search Bar -->
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
                    <!-- Search Results Overlay -->
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
                                    <span class="text-caption text-medium-emphasis">{{ formatMarket(item.market) }}</span>
                                </v-list-item-subtitle>
                                <template #append>
                                    <div class="d-flex ga-1">
                                        <v-btn
                                            v-if="!isInWatchlist(item.id)"
                                            size="x-small"
                                            variant="tonal"
                                            color="primary"
                                            @mousedown.prevent.stop="addToWatchlist(item)"
                                        >
                                            +{{ tt('Watchlist') }}
                                        </v-btn>
                                        <v-btn
                                            size="x-small"
                                            variant="tonal"
                                            color="success"
                                            @mousedown.prevent.stop="onBuyClick(item)"
                                        >
                                            {{ tt('Buy') }}
                                        </v-btn>
                                    </div>
                                </template>
                            </v-list-item>
                        </v-list>
                    </v-card>
                </div>
                <v-btn
                    v-if="showManageButton"
                    variant="outlined"
                    color="primary"
                    class="manage-btn ml-3"
                    @click="onManageClick"
                >
                    {{ tt('asset.Global Asset Management') }}
                </v-btn>
            </div>

            <!-- Tabs + Table (merged card) -->
            <v-card>
                <v-tabs v-model="activeTab" class="px-4">
                    <v-tab value="holdings">
                        {{ tt('Holdings') }}
                        <v-chip size="small" variant="tonal" class="ml-2">{{ filteredHoldings.length }}</v-chip>
                    </v-tab>
                    <v-tab value="watchlist">
                        {{ tt('Watchlist') }}
                        <v-chip size="small" variant="tonal" class="ml-2">{{ filteredWatchlist.length }}</v-chip>
                    </v-tab>
                </v-tabs>

                <!-- Holdings Tab Content -->
                <v-card-text v-if="activeTab === 'holdings'" class="pa-0">
                    <v-data-table
                        :headers="holdingsHeaders"
                        :items="filteredHoldings"
                        :loading="loading"
                        :sort-by="holdingsSortBy"
                        :hover="true"
                        item-value="assetId"
                        class="holdings-table"
                        v-model:items-per-page="holdingsPerPage"
                        v-model:page="holdingsPage"
                        @click:row="onRowClick"
                    >
                        <template #item.assetCode="{ item }">
                            <span class="text-body-2 font-weight-medium">{{ item.assetCode }}</span>
                        </template>
                        <template #item.assetName="{ item }">
                            <span class="text-body-2">{{ item.assetName }}</span>
                        </template>
                        <template #item.market="{ item }">
                            <span class="text-body-2">{{ formatMarket(item.market) }}</span>
                        </template>
                        <template #item.industry="{ item }">
                            <span v-if="item.industry" class="text-body-2">{{ formatIndustry(item.industry) }}</span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.quantity="{ item }">
                            <span v-if="item.quantity !== undefined && item.quantity !== null" class="text-body-2">
                                {{ formatQuantity(item.quantity) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.currentPrice="{ item }">
                            <span v-if="item.currentPrice !== undefined && item.currentPrice !== null" class="text-body-2">
                                {{ formatPrice(item.currentPrice) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.marketValue="{ item }">
                            <span v-if="item.marketValue !== undefined && item.marketValue !== null" class="text-body-2 font-weight-medium">
                                {{ formatCurrencyValue(item.marketValue, item.currency) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.totalCost="{ item }">
                            <span v-if="item.totalCost !== undefined && item.totalCost !== null" class="text-body-2">
                                {{ formatCurrencyValue(item.totalCost, item.currency) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.unrealizedPnl="{ item }">
                            <span v-if="item.unrealizedPnl !== undefined && item.unrealizedPnl !== null" class="text-body-2" :class="getReturnColorClass(item.unrealizedPnl)">
                                {{ formatCurrencyValue(item.unrealizedPnl, item.currency) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.returnRate="{ item }">
                            <span v-if="item.returnRate !== undefined && item.returnRate !== null" class="text-body-2" :class="getReturnColorClass(item.returnRate)">
                                {{ formatReturnRate(item.returnRate) }}
                            </span>
                            <span v-else class="text-medium-emphasis">--</span>
                        </template>
                        <template #item.actions="{ item }">
                            <div class="d-flex ga-1">
                                <v-btn size="x-small" variant="tonal" color="success" @click.stop="onBuyClick(item)">
                                    {{ tt('Buy') }}
                                </v-btn>
                                <v-btn size="x-small" variant="tonal" color="error" @click.stop="onSellClick(item)">
                                    {{ tt('Sell') }}
                                </v-btn>
                                <v-btn size="x-small" variant="text" icon="mdi-dots-horizontal" @click.stop="showDetail(item)" />
                            </div>
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
                                <pagination-buttons :totalPageCount="holdingsTotalPageCount"
                                                    v-model="holdingsPage"></pagination-buttons>
                            </div>
                        </template>
                    </v-data-table>
                </v-card-text>

                <!-- Watchlist Tab Content -->
                <v-card-text v-if="activeTab === 'watchlist'" class="pa-0">
                    <v-data-table
                        :headers="watchlistHeaders"
                        :items="filteredWatchlist"
                        :loading="loading"
                        :sort-by="watchlistSortBy"
                        :hover="true"
                        item-value="assetId"
                        class="watchlist-table"
                        v-model:items-per-page="watchlistPerPage"
                        v-model:page="watchlistPage"
                        @click:row="onRowClick"
                    >
                        <template #item.assetCode="{ item }">
                            <span class="text-body-2 font-weight-medium">{{ item.assetCode }}</span>
                        </template>
                        <template #item.assetName="{ item }">
                            <span class="text-body-2">{{ item.assetName }}</span>
                        </template>
                        <template #item.market="{ item }">
                            <span class="text-body-2">{{ formatMarket(item.market) }}</span>
                        </template>
                        <template #item.industry="{ item }">
                            <span v-if="item.industry" class="text-body-2">{{ formatIndustry(item.industry) }}</span>
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
                                <v-btn size="x-small" variant="tonal" color="success" @click.stop="onBuyClick(item)">
                                    {{ tt('Buy') }}
                                </v-btn>
                                <v-btn size="x-small" variant="tonal" color="error" @click.stop="removeFromWatchlist(item)">
                                    {{ tt('Remove') }}
                                </v-btn>
                                <v-btn size="x-small" variant="text" icon="mdi-dots-horizontal" @click.stop="showDetail(item)" />
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
                                <pagination-buttons :totalPageCount="watchlistTotalPageCount"
                                                    v-model="watchlistPage"></pagination-buttons>
                            </div>
                        </template>
                    </v-data-table>
                </v-card-text>

                <!-- Holdings Summary Footer (only under holdings tab) -->
                <template v-if="activeTab === 'holdings'">
                    <v-divider />
                    <v-card-text class="py-3 px-6">
                        <div class="d-flex flex-wrap align-center ga-6">
                            <div class="summary-item">
                                <span class="text-caption text-medium-emphasis">{{ tt('asset.TotalMarketValue') }}</span>
                                <span class="text-body-1 font-weight-bold ml-2">{{ formatCurrencyValue(holdingsSummary.totalMarketValue, 'CNY') }}</span>
                            </div>
                            <div class="summary-item">
                                <span class="text-caption text-medium-emphasis">{{ tt('Total Cost') }}</span>
                                <span class="text-body-1 font-weight-bold ml-2">{{ formatCurrencyValue(holdingsSummary.totalCost, 'CNY') }}</span>
                            </div>
                            <div class="summary-item">
                                <span class="text-caption text-medium-emphasis">{{ tt('Unrealized P&L') }}</span>
                                <span class="text-body-1 font-weight-bold ml-2" :class="getReturnColorClass(holdingsSummary.totalUnrealizedPnl)">
                                    {{ formatCurrencyValue(holdingsSummary.totalUnrealizedPnl, 'CNY') }}
                                </span>
                            </div>
                            <div class="summary-item">
                                <span class="text-caption text-medium-emphasis">{{ tt('Return Rate') }}</span>
                                <span class="text-body-1 font-weight-bold ml-2" :class="getReturnColorClass(holdingsSummary.totalReturnRate)">
                                    {{ formatReturnRate(holdingsSummary.totalReturnRate) }}
                                </span>
                            </div>
                        </div>
                    </v-card-text>
                </template>
            </v-card>

            <!-- Detail Dialog -->
            <v-dialog v-model="detailDialog" max-width="520">
                <v-card v-if="selectedItem">
                    <v-card-title class="d-flex align-center">
                        <span class="text-h6">{{ selectedItem.assetName }}</span>
                        <v-spacer />
                        <v-chip size="small" variant="tonal" :color="getCategoryColor(selectedItem.category)">
                            {{ formatCategory(selectedItem.category) }}
                        </v-chip>
                    </v-card-title>
                    <v-card-text>
                        <v-row dense>
                            <v-col cols="6">
                                <div class="text-caption text-medium-emphasis">{{ tt('Code') }}</div>
                                <div class="text-body-1">{{ selectedItem.assetCode }}</div>
                            </v-col>
                            <v-col cols="6">
                                <div class="text-caption text-medium-emphasis">{{ tt('Market') }}</div>
                                <div class="text-body-1">{{ formatMarket(selectedItem.market) }}</div>
                            </v-col>
                            <v-col cols="6">
                                <div class="text-caption text-medium-emphasis">{{ tt('Category') }}</div>
                                <div class="text-body-1">{{ formatCategory(selectedItem.category) }}</div>
                            </v-col>
                            <v-col cols="6">
                                <div class="text-caption text-medium-emphasis">{{ tt('Currency') }}</div>
                                <div class="text-body-1">{{ selectedItem.currency }}</div>
                            </v-col>
                        </v-row>

                        <template v-if="selectedItem.isHolding">
                            <v-divider class="my-3" />
                            <v-row dense>
                                <v-col cols="6">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Current Price') }}</div>
                                    <div class="text-body-1 font-weight-medium">
                                        {{ formatPrice(selectedItem.currentPrice) }}
                                    </div>
                                </v-col>
                                <v-col cols="6">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Quantity') }}</div>
                                    <div class="text-body-1">{{ formatQuantity(selectedItem.quantity) }}</div>
                                </v-col>
                                <v-col cols="6">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Avg Cost') }}</div>
                                    <div class="text-body-1">{{ formatPrice(selectedItem.avgCostPrice) }}</div>
                                </v-col>
                                <v-col cols="6">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Total Cost') }}</div>
                                    <div class="text-body-1">{{ formatCurrencyValue(selectedItem.totalCost, selectedItem.currency) }}</div>
                                </v-col>
                                <v-col cols="6">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Market Value') }}</div>
                                    <div class="text-body-1 font-weight-medium">{{ formatCurrencyValue(selectedItem.marketValue, selectedItem.currency) }}</div>
                                </v-col>
                                <v-col cols="6">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Unrealized P&L') }}</div>
                                    <div class="text-body-1" :class="getReturnColorClass(selectedItem.unrealizedPnl)">
                                        {{ formatCurrencyValue(selectedItem.unrealizedPnl, selectedItem.currency) }}
                                    </div>
                                </v-col>
                                <v-col cols="12">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Return Rate') }}</div>
                                    <div class="text-body-1 font-weight-bold" :class="getReturnColorClass(selectedItem.returnRate)">
                                        {{ formatReturnRate(selectedItem.returnRate) }}
                                    </div>
                                </v-col>
                            </v-row>
                        </template>

                        <template v-else>
                            <v-divider class="my-3" />
                            <v-row dense>
                                <v-col cols="12">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Current Price') }}</div>
                                    <div class="text-body-1" v-if="selectedItem.currentPrice !== undefined && selectedItem.currentPrice !== null">
                                        {{ formatPrice(selectedItem.currentPrice) }}
                                    </div>
                                    <div class="text-body-1 text-medium-emphasis" v-else>--</div>
                                </v-col>
                            </v-row>
                        </template>
                    </v-card-text>
                    <v-card-actions>
                        <v-spacer />
                        <v-btn variant="text" @click="detailDialog = false">{{ tt('Close') }}</v-btn>
                        <template v-if="!selectedItem.isHolding">
                            <v-btn color="error" variant="text" @click="removeFromWatchlist(selectedItem)">
                                {{ tt('Remove') }}
                            </v-btn>
                        </template>
                    </v-card-actions>
                </v-card>
            </v-dialog>

            <!-- Buy/Sell Dialog -->
            <v-dialog v-model="transactionDialog" max-width="520" persistent>
                <v-card v-if="selectedTransactionAsset">
                    <v-card-title class="d-flex align-center">
                        <span class="text-h6">{{ selectedTransactionAsset.name || selectedTransactionAsset.code }}</span>
                        <v-spacer />
                        <span class="text-body-2 text-medium-emphasis">{{ selectedTransactionAsset.code }}</span>
                    </v-card-title>
                    <v-card-text>
                        <v-form ref="transactionFormRef" @submit.prevent="submitTransaction">
                            <v-row dense>
                                <v-col cols="12">
                                    <v-select
                                        v-model="transactionForm.accountId"
                                        :label="tt('Account')"
                                        :items="investmentAccountItems"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        required
                                    />
                                </v-col>
                                <v-col cols="12">
                                    <v-select
                                        v-model="transactionForm.type"
                                        :label="tt('Transaction Type')"
                                        :items="transactionTypeOptions"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        required
                                    />
                                </v-col>
                                <v-col cols="12" md="6">
                                    <v-text-field
                                        v-model.number="transactionForm.amount"
                                        :label="tt('Amount')"
                                        type="number"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        step="0.01"
                                        min="0"
                                    />
                                </v-col>
                                <v-col cols="12" md="6">
                                    <v-text-field
                                        v-model.number="transactionForm.quantity"
                                        :label="tt('Quantity')"
                                        type="number"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        step="0.01"
                                        min="0"
                                    />
                                </v-col>
                                <v-col cols="12" md="6">
                                    <v-text-field
                                        v-model.number="transactionForm.price"
                                        :label="tt('Current Price')"
                                        type="number"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        step="0.0001"
                                        min="0"
                                    />
                                </v-col>
                                <v-col cols="12" md="6">
                                    <v-text-field
                                        v-model.number="transactionForm.fee"
                                        :label="tt('Service Charge')"
                                        type="number"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        step="0.01"
                                        min="0"
                                    />
                                </v-col>
                                <v-col cols="12">
                                    <v-text-field
                                        v-model="transactionForm.tradeTime"
                                        :label="tt('Transaction Time')"
                                        type="datetime-local"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                    />
                                </v-col>
                                <v-col cols="12">
                                    <v-textarea
                                        v-model="transactionForm.comment"
                                        :label="tt('Comment')"
                                        density="compact"
                                        variant="outlined"
                                        hide-details
                                        rows="2"
                                        :placeholder="tt('Comment')"
                                    />
                                </v-col>
                            </v-row>
                        </v-form>
                    </v-card-text>
                    <v-card-actions>
                        <v-spacer />
                        <v-btn variant="text" :disabled="transactionSubmitting" @click="transactionDialog = false">
                            {{ tt('Cancel') }}
                        </v-btn>
                        <v-btn
                            color="primary"
                            variant="elevated"
                            :disabled="!isTransactionFormValid || transactionSubmitting"
                            :loading="transactionSubmitting"
                            @click="submitTransaction"
                        >
                            {{ tt('Save') }}
                        </v-btn>
                    </v-card-actions>
                </v-card>
            </v-dialog>

            <!-- Admin Dialog -->
            <v-dialog v-model="adminDialog" max-width="900" scrollable>
                <v-card>
                    <v-card-title class="d-flex align-center">
                        <span class="text-h6">{{ tt('asset.Global Asset Management') }}</span>
                        <v-spacer />
                        <v-btn variant="text" icon="mdi-close" density="compact" @click="adminDialog = false" />
                    </v-card-title>
                    <v-divider />
                    <v-card-text class="pa-0">
                        <div class="pa-4 pb-2">
                            <div class="d-flex ga-3 align-center">
                                <v-text-field
                                    v-model="adminSearchKeyword"
                                    :prepend-inner-icon="mdiMagnify"
                                    density="compact"
                                    :placeholder="tt('asset.SearchAssetsByNameOrCode')"
                                    hide-details
                                    clearable
                                    variant="outlined"
                                    style="max-width: 360px"
                                    @update:model-value="onAdminSearchChange"
                                />
                                <v-select
                                    v-model="adminIndustryFilter"
                                    density="compact"
                                    :placeholder="tt('Industry')"
                                    :items="adminIndustryOptions"
                                    item-title="label"
                                    item-value="value"
                                    clearable
                                    hide-details
                                    variant="outlined"
                                    style="max-width: 180px"
                                />
                                <v-spacer />
                                <v-btn color="primary" variant="tonal" @click="onAssetAddClick">
                                    {{ tt('Add') }}
                                </v-btn>
                                <v-btn
                                    variant="tonal"
                                    color="secondary"
                                    :loading="syncingAssets"
                                    :disabled="syncingAssets"
                                    @click="onSyncAssets"
                                >
                                    <v-icon start size="small">mdi-database-sync</v-icon>
                                    {{ tt('asset.Sync Assets') }}
                                </v-btn>
                            </div>
                        </div>

                        <v-data-table
                            :headers="adminHeaders"
                            :items="filteredAdminAssets"
                            :loading="adminLoading"
                            :hover="true"
                            item-value="id"
                            density="compact"
                            class="admin-table"
                        >
                            <template #item.code="{ item }">
                                <span class="text-body-2 font-weight-medium">{{ item.code }}</span>
                            </template>
                            <template #item.name="{ item }">
                                <span class="text-body-2">{{ item.name }}</span>
                            </template>
                            <template #item.market="{ item }">
                                <span class="text-body-2">{{ formatMarket(item.market) }}</span>
                            </template>
                            <template #item.industry="{ item }">
                                <span v-if="item.industry" class="text-body-2">{{ formatIndustry(item.industry) }}</span>
                                <span v-else class="text-medium-emphasis">--</span>
                            </template>
                            <template #item.subCategory="{ item }">
                                <span v-if="item.subCategory" class="text-body-2">{{ item.subCategory }}</span>
                                <span v-else class="text-medium-emphasis">--</span>
                            </template>
                            <template #item.actions="{ item }">
                                <div class="d-flex ga-1">
                                    <v-btn size="x-small" variant="text" icon="mdi-pencil" @click.stop="onAssetEditClick(item)" />
                                    <v-btn size="x-small" variant="text" icon="mdi-delete" color="error" @click.stop="onAssetDeleteClick(item)" />
                                </div>
                            </template>
                            <template #loading>
                                <v-skeleton-loader type="table-row@10" :loading="true" />
                            </template>
                            <template #no-data>
                                <div class="text-center py-6 text-medium-emphasis">{{ tt('asset.NoAssetsFound') }}</div>
                            </template>
                        </v-data-table>
                    </v-card-text>
                    <v-divider />
                    <v-card-actions>
                        <v-spacer />
                        <v-btn variant="text" @click="adminDialog = false">{{ tt('Close') }}</v-btn>
                    </v-card-actions>
                </v-card>
            </v-dialog>

            <!-- Asset Create/Edit Dialog -->
            <v-dialog v-model="assetFormDialog" max-width="500">
                <v-card>
                    <v-card-title class="d-flex align-center">
                        <span class="text-h6">{{ assetFormMode === 'create' ? tt('Add') : tt('Edit') }}</span>
                        <v-spacer />
                        <v-btn variant="text" icon="mdi-close" density="compact" @click="assetFormDialog = false" />
                    </v-card-title>
                    <v-divider />
                    <v-card-text class="pa-4">
                        <v-text-field v-model="assetForm.code" :label="tt('asset.AssetCode')" density="compact" variant="outlined" :rules="[v => !!v || 'Required']" />
                        <v-text-field v-model="assetForm.name" :label="tt('asset.AssetName')" density="compact" variant="outlined" :rules="[v => !!v || 'Required']" class="mt-3" />
                        <div class="d-flex ga-3 mt-3">
                            <v-select v-model="assetForm.market" :label="tt('Market')" density="compact" variant="outlined" :items="marketOptions" style="flex: 1" />
                            <v-select v-model="assetForm.category" :label="tt('Category')" density="compact" variant="outlined" :items="categoryOptions" style="flex: 1" />
                        </div>
                        <div class="d-flex ga-3 mt-3">
                            <v-text-field v-model="assetForm.currency" :label="tt('Currency')" density="compact" variant="outlined" style="flex: 1" />
                            <v-select v-model="assetForm.industry" :label="tt('Industry')" density="compact" variant="outlined" :items="industryFormOptions" clearable style="flex: 1" />
                        </div>
                    </v-card-text>
                    <v-divider />
                    <v-card-actions>
                        <v-spacer />
                        <v-btn variant="text" @click="assetFormDialog = false">{{ tt('Cancel') }}</v-btn>
                        <v-btn color="primary" variant="tonal" :loading="assetFormSaving" :disabled="!isAssetFormValid" @click="onAssetFormSubmit">{{ tt('Save') }}</v-btn>
                    </v-card-actions>
                </v-card>
            </v-dialog>

            <!-- Asset Delete Confirm Dialog -->
            <v-dialog v-model="assetDeleteDialog" max-width="400">
                <v-card>
                    <v-card-title>{{ tt('Confirm') }}</v-card-title>
                    <v-card-text>
                        <span class="text-body-2">{{ tt('asset.ConfirmDeleteAsset') }} "{{ assetToDelete?.name }}" ?</span>
                    </v-card-text>
                    <v-card-actions>
                        <v-spacer />
                        <v-btn variant="text" @click="assetDeleteDialog = false">{{ tt('Cancel') }}</v-btn>
                        <v-btn color="error" variant="tonal" :loading="assetDeleteSaving" @click="onAssetDeleteConfirm">{{ tt('Delete') }}</v-btn>
                    </v-card-actions>
                </v-card>
            </v-dialog>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';

import PaginationButtons from '@/components/desktop/PaginationButtons.vue';

import { useI18n } from '@/locales/helpers.ts';

import { useInvestmentStore } from '@/stores/investment.ts';
import { useAccountsStore } from '@/stores/account.ts';

import {
    AssetCategory,
    InvestmentMarket,
    InvestmentTransactionType,
    InvestmentHolding,
    InvestmentUserAsset,
    InvestmentTransactionItem,
    type AssetInfoResponse,
    type AssetModifyRequest
} from '@/models/investment.ts';

import { mdiMagnify } from '@mdi/js';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';
import { getCurrentUnixTime, getTimezoneOffsetMinutes } from '@/lib/datetime.ts';

const { tt } = useI18n();

const investmentStore = useInvestmentStore();
const accountsStore = useAccountsStore();

const DIVISOR = 10000;

// --- State ---
const activeTab = ref<string>('holdings');
const searchQuery = ref<string>('');
const loading = ref<boolean>(true);
const detailDialog = ref<boolean>(false);
const selectedItem = ref<DisplayAsset | null>(null);
const showManageButton = ref<boolean>(false);
const adminDialog = ref<boolean>(false);
const adminLoading = ref<boolean>(false);
const adminSearchKeyword = ref<string>('');
const adminIndustryFilter = ref<string>('');
const adminAssets = ref<AssetInfoResponse[]>([]);
const assetFormDialog = ref<boolean>(false);
const assetFormMode = ref<'create' | 'edit'>('create');
const assetFormSaving = ref<boolean>(false);
const assetForm = ref({
    id: '',
    code: '',
    name: '',
    market: InvestmentMarket.CN as number,
    category: AssetCategory.Equity as string,
    currency: 'CNY',
    industry: ''
});
const assetDeleteDialog = ref<boolean>(false);
const assetDeleteSaving = ref<boolean>(false);
const assetToDelete = ref<AssetInfoResponse | null>(null);
const syncingAssets = ref<boolean>(false);
const syncResult = ref<{ fundsAdded: number; stocksAdded: number; etfsAdded: number; totalAdded: number } | null>(null);

const adminHeaders = [
    { key: 'code', title: tt('asset.AssetCode'), sortable: false },
    { key: 'name', title: tt('asset.AssetName'), sortable: false },
    { key: 'market', title: tt('Market'), sortable: false },
    { key: 'category', title: tt('Category'), sortable: false },
    { key: 'industry', title: tt('Industry') + ' ▾', sortable: false },
    { key: 'subCategory', title: tt('asset.SubCategory'), sortable: false },
    { key: 'currency', title: tt('Currency'), sortable: false },
    { key: 'actions', title: '', sortable: false, width: '100' }
];

const adminIndustryOptions = computed(() => {
    const industries = new Set<string>();
    adminAssets.value.forEach(a => { if (a.industry) industries.add(a.industry); });
    return Array.from(industries).sort().map(v => ({ label: formatIndustry(v), value: v }));
});

const filteredAdminAssets = computed(() => {
    let result = adminAssets.value;
    if (adminIndustryFilter.value) {
        result = result.filter(a => a.industry === adminIndustryFilter.value);
    }
    return result;
});

const marketOptions = [
    { title: 'A股', value: InvestmentMarket.CN },
    { title: '港股', value: InvestmentMarket.HK },
    { title: '美股', value: InvestmentMarket.US }
];

const categoryOptions = [
    { title: tt('Equity'), value: AssetCategory.Equity },
    { title: tt('Fixed Income'), value: AssetCategory.FixedIncome },
    { title: tt('Commodity'), value: AssetCategory.Commodity },
    { title: tt('Digital'), value: AssetCategory.Digital }
];

const industryFormOptions = ['technology', 'healthcare', 'consumer', 'finance', 'energy', 'industrial', 'real_estate', 'materials', 'utilities', 'telecom'].map(v => ({
    title: formatIndustry(v),
    value: v
}));

const isAssetFormValid = computed(() => !!assetForm.value.code && !!assetForm.value.name);

// --- Search State ---
const searchResults = ref<AssetInfoResponse[]>([]);
const searchLoading = ref<boolean>(false);
const showSearchResults = ref<boolean>(false);
let searchTimer: ReturnType<typeof setTimeout> | null = null;

// --- Watchlist ---
const watchlistAssets = ref<DisplayAsset[]>([]);
const watchlistAssetIds = ref<Set<string>>(new Set());

// --- Transaction Dialog ---
const transactionDialog = ref<boolean>(false);
const transactionSubmitting = ref<boolean>(false);
const selectedTransactionAsset = ref<{ assetId: string; code: string; name: string } | null>(null);
const transactionForm = ref({
    accountId: '',
    type: InvestmentTransactionType.Buy,
    amount: 0,
    quantity: 0,
    price: 0,
    fee: 0,
    tradeTime: '',
    comment: ''
});

// --- Pagination ---
const holdingsPerPage = ref<number>(10);
const holdingsPage = ref<number>(1);
const watchlistPerPage = ref<number>(10);
const watchlistPage = ref<number>(1);

const holdingsTotalPageCount = computed<number>(() => {
    return Math.max(1, Math.ceil(filteredHoldings.value.length / (holdingsPerPage.value > 0 ? holdingsPerPage.value : 10)));
});

const watchlistTotalPageCount = computed<number>(() => {
    return Math.max(1, Math.ceil(filteredWatchlist.value.length / (watchlistPerPage.value > 0 ? watchlistPerPage.value : 10)));
});

// --- Display Asset interface ---
interface DisplayAsset {
    assetId: string;
    assetCode: string;
    assetName: string;
    category: string;
    currency: string;
    market: number;
    industry: string;
    quantity?: number;
    avgCostPrice?: number;
    totalCost?: number;
    currentPrice?: number;
    marketValue?: number;
    unrealizedPnl?: number;
    returnRate?: number;
    isHolding: boolean;
}

// --- Transaction helpers ---
const investmentAccountItems = computed(() => {
    return accountsStore.allAccounts
        .filter(a => a.category === 7)
        .map(a => ({
            title: a.name,
            value: a.id
        }));
});

const transactionTypeOptions = computed(() => [
    { title: tt('Buy'), value: InvestmentTransactionType.Buy },
    { title: tt('Sell'), value: InvestmentTransactionType.Sell }
]);

const isTransactionFormValid = computed(() => {
    return transactionForm.value.accountId !== ''
        && transactionForm.value.amount > 0
        && transactionForm.value.tradeTime !== '';
});

// --- Table headers ---
const holdingsHeaders = computed(() => [
    { key: 'assetCode', title: tt('Code'), sortable: false },
    { key: 'assetName', title: tt('asset.AssetName'), sortable: false },
    { key: 'market', title: tt('Market') + ' ▾', sortable: false },
    { key: 'industry', title: tt('Industry') + ' ▾', sortable: false },
    { key: 'quantity', title: tt('Quantity'), sortable: false, align: 'end' as const },
    { key: 'currentPrice', title: tt('Current Price'), sortable: false, align: 'end' as const },
    { key: 'marketValue', title: tt('Market Value'), sortable: false, align: 'end' as const },
    { key: 'totalCost', title: tt('Total Cost'), sortable: false, align: 'end' as const },
    { key: 'unrealizedPnl', title: tt('Unrealized P&L'), sortable: false, align: 'end' as const },
    { key: 'returnRate', title: tt('Return Rate'), sortable: false, align: 'end' as const },
    { key: 'actions', title: tt('Action'), sortable: false, align: 'center' as const, width: '160' }
]);

const watchlistHeaders = computed(() => [
    { key: 'assetCode', title: tt('Code'), sortable: false },
    { key: 'assetName', title: tt('asset.AssetName'), sortable: false },
    { key: 'market', title: tt('Market') + ' ▾', sortable: false },
    { key: 'industry', title: tt('Industry') + ' ▾', sortable: false },
    { key: 'currentPrice', title: tt('Current Price'), sortable: false, align: 'end' as const },
    { key: 'actions', title: tt('Action'), sortable: false, align: 'center' as const, width: '160' }
]);

const holdingsSortBy = computed(() => []);

const watchlistSortBy = computed(() => []);

// --- Holdings ---
const holdingsList = computed<DisplayAsset[]>(() => {
    return investmentStore.holdings.map(h => holdingToDisplay(h));
});

function holdingToDisplay(h: InvestmentHolding): DisplayAsset {
    return {
        assetId: h.assetId,
        assetCode: h.assetCode,
        assetName: h.assetName,
        category: h.category,
        currency: h.currency,
        market: h.market,
        industry: '',
        quantity: h.quantity,
        avgCostPrice: h.avgCostPrice,
        totalCost: h.totalCost,
        currentPrice: h.currentPrice,
        marketValue: h.marketValue,
        unrealizedPnl: h.unrealizedPnl,
        returnRate: h.returnRate,
        isHolding: true
    };
}

// --- Holdings Summary ---
const holdingsSummary = computed(() => {
    const holdings = filteredHoldings.value;
    let totalMarketValue = 0;
    let totalUnrealizedPnl = 0;
    let totalCost = 0;

    for (const h of holdings) {
        if (h.marketValue !== undefined && h.marketValue !== null) {
            totalMarketValue += h.marketValue;
        }
        if (h.unrealizedPnl !== undefined && h.unrealizedPnl !== null) {
            totalUnrealizedPnl += h.unrealizedPnl;
        }
        if (h.totalCost !== undefined && h.totalCost !== null) {
            totalCost += h.totalCost;
        }
    }

    const totalReturnRate = totalCost !== 0
        ? Math.round((totalUnrealizedPnl / totalCost) * 10000)
        : 0;

    return {
        totalMarketValue,
        totalCost,
        totalUnrealizedPnl,
        totalReturnRate
    };
});

// --- Filtered lists ---
const filteredHoldings = computed<DisplayAsset[]>(() => {
    return holdingsList.value;
});

const filteredWatchlist = computed<DisplayAsset[]>(() => {
    return watchlistAssets.value;
});

// --- Search ---
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
    searchLoading.value = true;
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
    } finally {
        searchLoading.value = false;
    }
}

function onSearchResultClick(item: AssetInfoResponse): void {
    showSearchResults.value = false;
    searchQuery.value = '';

    const existingHolding = holdingsList.value.find(h => h.assetId === item.id);
    if (existingHolding) {
        showDetail(existingHolding);
        return;
    }

    const existingWatchlist = watchlistAssets.value.find(w => w.assetId === item.id);
    if (existingWatchlist) {
        showDetail(existingWatchlist);
        return;
    }

    const displayItem: DisplayAsset = {
        assetId: item.id,
        assetCode: item.code,
        assetName: item.name,
        category: item.category,
        currency: item.currency,
        market: item.market,
        industry: item.industry || '',
        isHolding: false
    };
    showDetail(displayItem);
}

function isInWatchlist(assetId: string): boolean {
    return watchlistAssetIds.value.has(assetId);
}

async function addToWatchlist(item: AssetInfoResponse): Promise<void> {
    try {
        await investmentStore.addUserAsset({ assetId: item.id });
        await loadWatchlist();
    } catch (error) {
        logger.error('Failed to add to watchlist', error);
    }
}

// --- Watchlist loading ---
async function loadWatchlist(): Promise<void> {
    try {
        const response = await services.getUserAssets({ is_watchlist: true });
        const data = response.data;

        if (!data || !data.success || !data.result) {
            logger.error('Failed to load watchlist assets');
            return;
        }

        const userAssets = InvestmentUserAsset.ofMulti(data.result);
        const holdingAssetIds = new Set(investmentStore.holdings.map(h => h.assetId));

        const watchlistItems: DisplayAsset[] = [];
        const ids: Set<string> = new Set();

        for (const ua of userAssets) {
            if (!ua.asset) {
                continue;
            }

            const a = ua.asset;
            ids.add(a.id);

            watchlistItems.push({
                assetId: a.id,
                assetCode: a.code,
                assetName: a.name,
                category: a.category,
                currency: a.currency,
                market: a.market,
                industry: a.industry || '',
                isHolding: holdingAssetIds.has(a.id)
            });
        }

        watchlistAssets.value = watchlistItems;
        watchlistAssetIds.value = ids;

        loadWatchlistPrices(watchlistItems);
    } catch (error) {
        logger.error('Failed to load watchlist', error);
    }
}

async function loadWatchlistPrices(items: DisplayAsset[]): Promise<void> {
    const pricePromises = items.map(async (item) => {
        try {
            const result = await investmentStore.loadLatestMarketData({ assetId: item.assetId });
            if (result) {
                item.currentPrice = result.price;
            }
        } catch {
            // Ignore individual failures
        }
    });

    await Promise.all(pricePromises);
}

async function removeFromWatchlist(item: DisplayAsset): Promise<void> {
    try {
        await investmentStore.removeUserAsset({ assetId: item.assetId });
        watchlistAssets.value = watchlistAssets.value.filter(a => a.assetId !== item.assetId);
        watchlistAssetIds.value.delete(item.assetId);
        detailDialog.value = false;
    } catch (error) {
        logger.error('Failed to remove asset from watchlist', error);
    }
}

// --- Detail dialog ---
function onRowClick(_event: Event, row: { item: DisplayAsset }): void {
    showDetail(row.item);
}

function showDetail(item: DisplayAsset): void {
    selectedItem.value = item;
    detailDialog.value = true;
}

// --- Action handlers ---
function onBuyClick(item: DisplayAsset | AssetInfoResponse): void {
    const isAssetInfo = (i: DisplayAsset | AssetInfoResponse): i is AssetInfoResponse => 'id' in i;
    if (isAssetInfo(item)) {
        selectedTransactionAsset.value = { assetId: item.id, code: item.code, name: item.name };
    } else {
        selectedTransactionAsset.value = { assetId: item.assetId, code: item.assetCode, name: item.assetName };
    }
    transactionForm.value.type = InvestmentTransactionType.Buy;
    transactionForm.value.amount = 0;
    transactionForm.value.quantity = 0;
    transactionForm.value.price = 0;
    transactionForm.value.fee = 0;
    transactionForm.value.tradeTime = '';
    transactionForm.value.comment = '';
    transactionDialog.value = true;
}

function onSellClick(item: DisplayAsset): void {
    selectedTransactionAsset.value = {
        assetId: item.assetId,
        code: item.assetCode,
        name: item.assetName
    };
    transactionForm.value.type = InvestmentTransactionType.Sell;
    transactionForm.value.amount = 0;
    transactionForm.value.quantity = 0;
    transactionForm.value.price = 0;
    transactionForm.value.fee = 0;
    transactionForm.value.tradeTime = '';
    transactionForm.value.comment = '';
    transactionDialog.value = true;
}

async function onManageClick(): Promise<void> {
    adminDialog.value = true;
    adminSearchKeyword.value = '';
    adminIndustryFilter.value = '';
    await loadAdminAssets();
}

async function onSyncAssets(): Promise<void> {
    syncingAssets.value = true;
    syncResult.value = null;
    try {
        const resp = await services.syncGlobalAssets();
        const data = resp.data?.result;
        if (data) {
            syncResult.value = data;
            await loadAdminAssets();
        }
    } catch (e) {
        logger.error('Failed to sync assets', e);
    } finally {
        syncingAssets.value = false;
    }
}

async function loadAdminAssets(): Promise<void> {
    adminLoading.value = true;
    try {
        const resp = await services.listGlobalAssets({
            keyword: adminSearchKeyword.value || undefined
        });
        adminAssets.value = resp.data?.result?.assets || [];
    } catch (e) {
        logger.error('Failed to load admin assets', e);
        adminAssets.value = [];
    } finally {
        adminLoading.value = false;
    }
}

function onAdminSearchChange(): void {
    loadAdminAssets();
}

function onAssetAddClick(): void {
    assetFormMode.value = 'create';
    assetForm.value = {
        id: '',
        code: '',
        name: '',
        market: InvestmentMarket.CN,
        category: AssetCategory.Equity,
        currency: 'CNY',
        industry: ''
    };
    assetFormDialog.value = true;
}

function onAssetEditClick(asset: AssetInfoResponse): void {
    assetFormMode.value = 'edit';
    assetForm.value = {
        id: asset.id,
        code: asset.code,
        name: asset.name,
        market: asset.market,
        category: asset.category,
        currency: asset.currency,
        industry: asset.industry || ''
    };
    assetFormDialog.value = true;
}

function onAssetDeleteClick(asset: AssetInfoResponse): void {
    assetToDelete.value = asset;
    assetDeleteDialog.value = true;
}

async function onAssetFormSubmit(): Promise<void> {
    if (!isAssetFormValid.value) return;
    assetFormSaving.value = true;
    try {
        if (assetFormMode.value === 'create') {
            await services.addGlobalAsset({
                code: assetForm.value.code,
                market: assetForm.value.market,
                name: assetForm.value.name,
                category: assetForm.value.category,
                currency: assetForm.value.currency,
                industry: assetForm.value.industry || undefined
            });
        } else {
            const req: AssetModifyRequest = {
                id: assetForm.value.id,
                code: assetForm.value.code,
                market: assetForm.value.market,
                name: assetForm.value.name,
                category: assetForm.value.category,
                currency: assetForm.value.currency,
                industry: assetForm.value.industry || undefined
            };
            await services.modifyGlobalAsset(req);
        }
        assetFormDialog.value = false;
        await loadAdminAssets();
    } catch (e) {
        logger.error('Failed to save asset', e);
    } finally {
        assetFormSaving.value = false;
    }
}

async function onAssetDeleteConfirm(): Promise<void> {
    if (!assetToDelete.value) return;
    assetDeleteSaving.value = true;
    try {
        await services.deleteGlobalAsset({ id: assetToDelete.value.id });
        assetDeleteDialog.value = false;
        assetToDelete.value = null;
        await loadAdminAssets();
    } catch (e) {
        logger.error('Failed to delete asset', e);
    } finally {
        assetDeleteSaving.value = false;
    }
}

// --- Submit Transaction ---
async function submitTransaction(): Promise<void> {
    if (!selectedTransactionAsset.value) return;

    const asset = selectedTransactionAsset.value;
    const form = transactionForm.value;

    transactionSubmitting.value = true;
    try {
        const accountId = form.accountId;
        const type = form.type;

        const tradeTimeStr = form.tradeTime;
        let tradeTime = getCurrentUnixTime();
        if (tradeTimeStr) {
            tradeTime = Math.floor(new Date(tradeTimeStr).getTime() / 1000);
        }

        const utcOffset = getTimezoneOffsetMinutes(tradeTime);

        const transaction = new InvestmentTransactionItem(
            '',
            asset.assetId,
            accountId,
            type,
            tradeTime,
            Math.round(form.amount * DIVISOR)
        );
        transaction.quantity = Math.round(form.quantity * DIVISOR);
        transaction.price = Math.round(form.price * DIVISOR);
        transaction.fee = Math.round(form.fee * DIVISOR);
        transaction.utcOffset = utcOffset;
        transaction.comment = form.comment;

        await investmentStore.addTransaction({ transaction });

        transactionDialog.value = false;
        investmentStore.holdingsStateInvalid = true;
        await investmentStore.loadHoldings({ force: true });
    } catch (error) {
        logger.error('Failed to submit transaction', error);
    } finally {
        transactionSubmitting.value = false;
    }
}

// --- Admin check ---
async function checkAdmin(): Promise<void> {
    try {
        const response = await services.checkInvestmentAdmin();
        const data = response.data;
        if (data && data.success && data.result) {
            showManageButton.value = data.result.isAdmin;
        }
    } catch (error) {
        logger.error('Failed to check admin status', error);
    }
}

// --- Format helpers ---
function formatPrice(value: number | undefined | null): string {
    if (value === undefined || value === null) return '--';
    const num = value / DIVISOR;
    if (num >= 1000) return num.toFixed(2);
    if (num >= 1) return num.toFixed(4);
    return num.toFixed(6);
}

function formatCurrencyValue(value: number | undefined | null, currency: string): string {
    if (value === undefined || value === null) return '--';
    const num = value / DIVISOR;
    const abs = Math.abs(num);
    let formatted: string;
    if (abs >= 1000000) {
        formatted = (abs / 1000000).toFixed(2) + 'M';
    } else if (abs >= 10000) {
        formatted = (abs / 10000).toFixed(2) + 'W';
    } else if (abs >= 1000) {
        formatted = (abs / 1000).toFixed(2) + 'K';
    } else {
        formatted = abs.toFixed(2);
    }
    const sign = num < 0 ? '-' : '';
    return sign + getCurrencySymbol(currency) + formatted;
}

function formatReturnRate(value: number | undefined | null): string {
    if (value === undefined || value === null) return '--';
    const pct = value / 100;
    const sign = pct >= 0 ? '+' : '';
    return sign + pct.toFixed(2) + '%';
}

function formatQuantity(value: number | undefined | null): string {
    if (value === undefined || value === null) return '--';
    const num = value / DIVISOR;
    if (num >= 1000000) return (num / 1000000).toFixed(2) + 'M';
    if (num >= 10000) return (num / 10000).toFixed(2) + 'W';
    if (num >= 1000) return (num / 1000).toFixed(2) + 'K';
    return num.toFixed(2);
}

function formatMarket(market: number): string {
    switch (market) {
        case InvestmentMarket.CN: return tt('Market CN');
        case InvestmentMarket.HK: return tt('Market HK');
        case InvestmentMarket.US: return tt('Market US');
        default: return market.toString();
    }
}

function formatCategory(category: string): string {
    switch (category) {
        case AssetCategory.Equity: return tt('Equity');
        case AssetCategory.FixedIncome: return tt('Fixed Income');
        case AssetCategory.Commodity: return tt('Commodity');
        case AssetCategory.Digital: return tt('Digital');
        default: return category;
    }
}

function formatIndustry(industry: string): string {
    const industryMap: Record<string, string> = {
        'technology': tt('Technology'),
        'healthcare': tt('Healthcare'),
        'consumer': tt('Consumer'),
        'finance': tt('Finance'),
        'energy': tt('Energy'),
        'industrial': tt('Industrial'),
        'real_estate': tt('Real Estate'),
        'materials': tt('Materials'),
        'utilities': tt('Utilities'),
        'telecom': tt('Telecom')
    };
    return industryMap[industry] || industry;
}

function getCurrencySymbol(currency: string): string {
    switch (currency) {
        case 'CNY': return '\u00a5';
        case 'HKD': return 'HK$';
        case 'USD': return '$';
        default: return '';
    }
}

function getReturnColorClass(value: number | undefined | null): string {
    if (value === undefined || value === null) return '';
    if (value > 0) return 'text-success';
    if (value < 0) return 'text-error';
    return '';
}

function getCategoryColor(category: string): string {
    switch (category) {
        case AssetCategory.Equity: return 'primary';
        case AssetCategory.FixedIncome: return 'success';
        case AssetCategory.Commodity: return 'warning';
        case AssetCategory.Digital: return 'purple';
        default: return 'default';
    }
}

// --- Watch for tab change ---
watch(activeTab, () => {
    // Reset state on tab change if needed
});

// --- Lifecycle ---
onMounted(async () => {
    loading.value = true;
    try {
        checkAdmin();
        await accountsStore.loadAllAccounts({ force: false });
        await investmentStore.loadHoldings({ force: false });
        await loadWatchlist();
    } catch (error) {
        logger.error('Failed to load assets page data', error);
    } finally {
        loading.value = false;
    }
});
</script>

<style scoped>
.page-content {
    padding: 24px;
}

.page-body {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

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
    max-height: 400px;
    overflow-y: auto;
}

.summary-item {
    display: flex;
    align-items: center;
    white-space: nowrap;
}

.holdings-table :deep(.v-data-table__td),
.watchlist-table :deep(.v-data-table__td) {
    padding-top: 8px;
    padding-bottom: 8px;
}
</style>

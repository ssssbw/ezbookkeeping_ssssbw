<template>
    <div class="page-content">
        <div class="page-header">
            <div class="d-flex align-center">
                <h1 class="page-title">{{ tt('Investment Transactions') }}</h1>
                <v-btn color="primary" variant="tonal" class="ms-4" @click="showCreateDialog = true">
                    <v-icon :icon="mdiPlus" class="me-1" />
                    {{ tt('Add Transaction') }}
                </v-btn>
            </div>
        </div>
        <div class="page-body">
            <v-row>
                <v-col cols="12" md="4">
                    <v-card>
                        <v-card-text class="summary-card">
                            <div class="summary-label">{{ tt('Buy') }}</div>
                            <div class="summary-value text-primary">{{ formatCurrencyValue(summaryData.totalBuy, 'CNY') }}</div>
                            <div class="summary-count">{{ summaryData.buyCount }} {{ tt('transactions') }}</div>
                        </v-card-text>
                    </v-card>
                </v-col>
                <v-col cols="12" md="4">
                    <v-card>
                        <v-card-text class="summary-card">
                            <div class="summary-label">{{ tt('Sell') }}</div>
                            <div class="summary-value text-error">{{ formatCurrencyValue(summaryData.totalSell, 'CNY') }}</div>
                            <div class="summary-count">{{ summaryData.sellCount }} {{ tt('transactions') }}</div>
                        </v-card-text>
                    </v-card>
                </v-col>
                <v-col cols="12" md="4">
                    <v-card>
                        <v-card-text class="summary-card">
                            <div class="summary-label">{{ tt('Dividend') }}</div>
                            <div class="summary-value text-success">{{ formatCurrencyValue(summaryData.totalDividend, 'CNY') }}</div>
                            <div class="summary-count">{{ summaryData.dividendCount }} {{ tt('transactions') }}</div>
                        </v-card-text>
                    </v-card>
                </v-col>
            </v-row>

            <v-card>
                <v-card-text class="pa-0">
                    <v-data-table
                        :headers="headers"
                        :items="transactions"
                        :loading="loading"
                        :hover="true"
                        :items-per-page="15"
                        v-model:page="page"
                        class="transactions-table"
                    >
                        <template #item.tradeTime="{ item }">
                            <span class="text-body-2">{{ formatTradeTime(item.tradeTime) }}</span>
                        </template>
                        <template #item.type="{ item }">
                            <v-chip size="small" :color="getTypeColor(item.type)" variant="tonal">
                                {{ formatTransactionType(item.type, tt) }}
                            </v-chip>
                        </template>
                        <template #item.assetCode="{ item }">
                            <span class="text-body-2 asset-code">{{ item.assetCode }}</span>
                        </template>
                        <template #item.assetName="{ item }">
                            <span class="text-body-2">{{ item.assetName }}</span>
                        </template>
                        <template #item.quantity="{ item }">
                            <span class="text-body-2">{{ formatQuantity(item.quantity) }}</span>
                        </template>
                        <template #item.price="{ item }">
                            <span class="text-body-2">{{ formatPrice(item.price) }}</span>
                        </template>
                        <template #item.amount="{ item }">
                            <span class="text-body-2">{{ formatCurrencyValue(item.amount, 'CNY') }}</span>
                        </template>
                        <template #item.fee="{ item }">
                            <span class="text-body-2">{{ formatCurrencyValue(item.fee, 'CNY') }}</span>
                        </template>
                        <template #item.comment="{ item }">
                            <span class="text-body-2 text-medium-emphasis">{{ item.comment || '--' }}</span>
                        </template>
                        <template #loading>
                            <v-skeleton-loader type="table-row@10" :loading="true" />
                        </template>
                        <template #no-data>
                            <div class="text-center py-6 text-medium-emphasis">
                                {{ tt('No transaction data') }}
                            </div>
                        </template>
                        <template #bottom>
                            <div class="mt-2 mb-4">
                                <pagination-buttons :totalPageCount="totalPageCount" v-model="page" />
                            </div>
                        </template>
                    </v-data-table>
                </v-card-text>
            </v-card>
        </div>

        <TransactionFormDialog
            v-model="showCreateDialog"
            @save="onTransactionSaved"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { mdiPlus } from '@mdi/js';

import PaginationButtons from '@/components/desktop/PaginationButtons.vue';
import TransactionFormDialog from './components/TransactionFormDialog.vue';

import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';

import { InvestmentTransactionType } from '@/models/investment.ts';
import { formatPrice, formatQuantity, formatCurrencyValue } from './assets/assetUtils.ts';

const { tt } = useI18n();
const investmentStore = useInvestmentStore();

const loading = ref(true);
const page = ref(1);
const showCreateDialog = ref(false);

const transactions = computed(() => investmentStore.transactions);

const totalPageCount = computed(() => {
    return Math.max(1, Math.ceil(transactions.value.length / 15));
});

const headers = computed(() => [
    { key: 'tradeTime', title: tt('Date'), sortable: false, width: '110' },
    { key: 'type', title: tt('Type'), sortable: false, width: '100' },
    { key: 'assetCode', title: tt('Code'), sortable: false, width: '100' },
    { key: 'assetName', title: tt('asset.AssetName'), sortable: false },
    { key: 'quantity', title: tt('Quantity'), sortable: false },
    { key: 'price', title: tt('Price'), sortable: false },
    { key: 'amount', title: tt('Amount'), sortable: false },
    { key: 'fee', title: tt('Fee'), sortable: false },
    { key: 'comment', title: tt('Comment'), sortable: false },
]);

interface SummaryData {
    totalBuy: number;
    totalSell: number;
    totalDividend: number;
    buyCount: number;
    sellCount: number;
    dividendCount: number;
}

const summaryData = computed<SummaryData>(() => {
    const result: SummaryData = {
        totalBuy: 0, totalSell: 0, totalDividend: 0,
        buyCount: 0, sellCount: 0, dividendCount: 0,
    };
    for (const t of transactions.value) {
        switch (t.type) {
            case InvestmentTransactionType.Buy:
                result.totalBuy += t.amount;
                result.buyCount++;
                break;
            case InvestmentTransactionType.Sell:
                result.totalSell += t.amount;
                result.sellCount++;
                break;
            case InvestmentTransactionType.DividendCash:
            case InvestmentTransactionType.DividendReinvest:
                result.totalDividend += t.amount;
                result.dividendCount++;
                break;
        }
    }
    return result;
});

function formatTradeTime(timestamp: number): string {
    if (!timestamp) return '--';
    const d = new Date(timestamp * 1000);
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
}

function formatTransactionType(type: number, tt: (key: string) => string): string {
    switch (type) {
        case InvestmentTransactionType.Buy: return tt('Buy');
        case InvestmentTransactionType.Sell: return tt('Sell');
        case InvestmentTransactionType.DividendCash: return tt('Dividend');
        case InvestmentTransactionType.DividendReinvest: return tt('Dividend Reinvest');
        default: return type.toString();
    }
}

function getTypeColor(type: number): string {
    switch (type) {
        case InvestmentTransactionType.Buy: return 'primary';
        case InvestmentTransactionType.Sell: return 'error';
        case InvestmentTransactionType.DividendCash:
        case InvestmentTransactionType.DividendReinvest: return 'success';
        default: return 'default';
    }
}

async function loadTransactions(): Promise<void> {
    loading.value = true;
    try {
        await investmentStore.loadTransactions({ force: true });
    } finally {
        loading.value = false;
    }
}

function onTransactionSaved(): void {
    loadTransactions();
}

onMounted(() => {
    loadTransactions();
});
</script>

<style scoped>
.page-content {
    padding: 24px;
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow-y: auto;
}

.page-header {
    margin-bottom: 24px;
}

.page-title {
    font-size: 24px;
    font-weight: 600;
    margin: 0;
}

.page-body {
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.summary-card {
    text-align: center;
}

.summary-label {
    font-size: 12px;
    color: rgba(var(--v-theme-on-surface), 0.6);
    margin-bottom: 4px;
}

.summary-value {
    font-size: 20px;
    font-weight: 500;
}

.summary-count {
    font-size: 12px;
    color: rgba(var(--v-theme-on-surface), 0.5);
    margin-top: 2px;
}

.asset-code {
    color: rgb(var(--v-theme-primary));
    font-weight: 500;
}

.transactions-table :deep(.v-data-table__td) {
    padding-top: 8px;
    padding-bottom: 8px;
}
</style>

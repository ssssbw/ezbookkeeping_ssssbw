<template>
    <div class="page-content">
        <div class="page-header">
            <h1 class="page-title">{{ tt('Investment Transactions') }}</h1>
        </div>
        <div class="page-body">
            <v-card class="mb-4">
                <v-card-title>{{ tt('Transaction History') }}</v-card-title>
                <v-card-text>
                    <div class="overflow-x-auto">
                        <table class="w-full">
                            <thead>
                                <tr>
                                    <th class="text-left">Date</th>
                                    <th class="text-left">Type</th>
                                    <th class="text-left">Asset</th>
                                    <th class="text-right">Quantity</th>
                                    <th class="text-right">Price</th>
                                    <th class="text-right">Total</th>
                                    <th class="text-right">Status</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="(transaction, index) in transactions" :key="index">
                                    <td>{{ transaction.date }}</td>
                                    <td>{{ transaction.type }}</td>
                                    <td>{{ transaction.asset }}</td>
                                    <td class="text-right">{{ transaction.quantity }}</td>
                                    <td class="text-right">{{ transaction.price }}</td>
                                    <td class="text-right">{{ transaction.total }}</td>
                                    <td class="text-right">
                                        <span :class="transaction.status === 'Completed' ? 'text-green-600' : 'text-yellow-600'">{{ transaction.status }}</span>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </v-card-text>
                <v-card-actions>
                    <v-btn color="primary">{{ tt('Add Transaction') }}</v-btn>
                </v-card-actions>
            </v-card>
            <v-card class="mb-4">
                <v-card-title>{{ tt('Transaction Summary') }}</v-card-title>
                <v-card-text>
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <v-card variant="outlined" class="p-4">
                            <div class="text-center">
                                <div class="text-2xl font-bold">$15,000</div>
                                <div class="text-gray-500">Total Buys</div>
                            </div>
                        </v-card>
                        <v-card variant="outlined" class="p-4">
                            <div class="text-center">
                                <div class="text-2xl font-bold">$5,000</div>
                                <div class="text-gray-500">Total Sells</div>
                            </div>
                        </v-card>
                        <v-card variant="outlined" class="p-4">
                            <div class="text-center">
                                <div class="text-2xl font-bold">$100</div>
                                <div class="text-gray-500">Total Dividends</div>
                            </div>
                        </v-card>
                    </div>
                </v-card-text>
            </v-card>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from '@/locales/helpers.ts';

const { tt } = useI18n();

const transactions = ref([
    { date: '2026-03-15', type: 'Buy', asset: 'AAPL', quantity: '10', price: '$180.00', total: '$1,800', status: 'Completed' },
    { date: '2026-03-14', type: 'Sell', asset: 'MSFT', quantity: '5', price: '$400.00', total: '$2,000', status: 'Completed' },
    { date: '2026-03-13', type: 'Buy', asset: 'GOOGL', quantity: '10', price: '$130.00', total: '$1,300', status: 'Pending' },
    { date: '2026-03-12', type: 'Dividend', asset: 'AAPL', quantity: '', price: '$0.50', total: '$50', status: 'Completed' },
    { date: '2026-03-10', type: 'Buy', asset: 'SPY', quantity: '20', price: '$500.00', total: '$10,000', status: 'Completed' },
]);
</script>

<style scoped>
.page-content {
    padding: 24px;
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
</style>
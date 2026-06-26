<template>
    <v-dialog v-model="dialogVisible" max-width="600" persistent>
        <v-card>
            <v-card-title class="text-subtitle-1">{{ tt('Add Transaction') }}</v-card-title>
            <v-card-text>
                <v-form ref="formRef" v-model="valid">
                    <v-row>
                        <v-col cols="12">
                            <v-autocomplete
                                v-model="form.assetId"
                                :items="assetOptions"
                                :loading="assetsLoading"
                                item-title="name"
                                item-value="id"
                                :label="tt('asset.AssetCode')"
                                :rules="[v => !!v || 'Required']"
                                @update:search="onAssetSearch"
                                variant="outlined"
                                density="compact"
                                hide-no-data
                            >
                                <template #item="{ props, item }">
                                    <v-list-item v-bind="props">
                                        <template #prepend>
                                            <span class="asset-code">{{ item.raw.code }}</span>
                                        </template>
                                    </v-list-item>
                                </template>
                            </v-autocomplete>
                        </v-col>
                        <v-col cols="12" md="6">
                            <v-select
                                v-model="form.type"
                                :items="transactionTypes"
                                :label="tt('Type')"
                                :rules="[v => v !== null || 'Required']"
                                variant="outlined"
                                density="compact"
                            />
                        </v-col>
                        <v-col cols="12" md="6">
                            <v-text-field
                                v-model="form.tradeDate"
                                :label="tt('Trade Date')"
                                type="date"
                                :rules="[v => !!v || 'Required']"
                                variant="outlined"
                                density="compact"
                            />
                        </v-col>
                        <v-col cols="12" md="6">
                            <v-text-field
                                v-model.number="form.quantity"
                                :label="tt('Quantity')"
                                type="number"
                                step="0.01"
                                variant="outlined"
                                density="compact"
                            />
                        </v-col>
                        <v-col cols="12" md="6">
                            <v-text-field
                                v-model.number="form.price"
                                :label="tt('Price')"
                                type="number"
                                step="0.0001"
                                variant="outlined"
                                density="compact"
                            />
                        </v-col>
                        <v-col cols="12" md="6">
                            <v-text-field
                                v-model.number="form.amount"
                                :label="tt('Amount')"
                                type="number"
                                step="0.01"
                                :rules="[v => !!v || 'Required']"
                                variant="outlined"
                                density="compact"
                            />
                        </v-col>
                        <v-col cols="12" md="6">
                            <v-text-field
                                v-model.number="form.fee"
                                :label="tt('Fee')"
                                type="number"
                                step="0.01"
                                variant="outlined"
                                density="compact"
                            />
                        </v-col>
                        <v-col cols="12">
                            <v-text-field
                                v-model="form.comment"
                                :label="tt('Comment')"
                                variant="outlined"
                                density="compact"
                            />
                        </v-col>
                    </v-row>
                </v-form>
            </v-card-text>
            <v-card-actions>
                <v-spacer />
                <v-btn variant="text" @click="dialogVisible = false">{{ tt('Cancel') }}</v-btn>
                <v-btn color="primary" variant="tonal" :loading="saving" @click="save">{{ tt('Save') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';
import { InvestmentTransactionType } from '@/models/investment.ts';
import { InvestmentTransactionItem } from '@/models/investment.ts';
import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

interface AssetOption {
    id: string;
    code: string;
    name: string;
}

const props = defineProps<{
    modelValue: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void;
    (e: 'save'): void;
}>();

const { tt } = useI18n();
const investmentStore = useInvestmentStore();

const valid = ref(false);
const saving = ref(false);
const assetsLoading = ref(false);
const assetOptions = ref<AssetOption[]>([]);

const form = ref({
    assetId: '',
    type: InvestmentTransactionType.Buy as number | null,
    tradeDate: new Date().toISOString().split('T')[0],
    quantity: 0,
    price: 0,
    amount: 0,
    fee: 0,
    comment: '',
});

const transactionTypes = computed(() => [
    { title: tt('Buy'), value: InvestmentTransactionType.Buy },
    { title: tt('Sell'), value: InvestmentTransactionType.Sell },
    { title: tt('Dividend'), value: InvestmentTransactionType.DividendCash },
    { title: tt('Dividend Reinvest'), value: InvestmentTransactionType.DividendReinvest },
]);

const dialogVisible = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val),
});

async function onAssetSearch(query: string): Promise<void> {
    if (!query || query.length < 1) return;
    assetsLoading.value = true;
    try {
        const resp = await services.searchAssets({ keyword: query });
        const data = resp.data;
        if (data.success && data.result) {
            assetOptions.value = data.result.map((a: { id: string; code: string; name: string }) => ({
                id: a.id,
                code: a.code,
                name: `${a.code} - ${a.name}`,
            }));
        }
    } catch {
        assetOptions.value = [];
    } finally {
        assetsLoading.value = false;
    }
}

async function save(): Promise<void> {
    if (!valid.value || !form.value.assetId || form.value.type === null) return;

    saving.value = true;
    try {
        const tradeDate = new Date(form.value.tradeDate || Date.now());
        const tradeTime = Math.floor(tradeDate.getTime() / 1000);
        const utcOffset = -tradeDate.getTimezoneOffset();

        const item = new InvestmentTransactionItem(
            '',
            form.value.assetId,
            '',
            form.value.type,
            tradeTime,
            form.value.amount,
        );
        item.quantity = form.value.quantity;
        item.price = form.value.price;
        item.fee = form.value.fee;
        item.utcOffset = utcOffset;
        item.comment = form.value.comment;

        await investmentStore.addTransaction({ transaction: item });
        dialogVisible.value = false;
        resetForm();
        emit('save');
    } catch (error) {
        logger.error('Failed to create transaction', error);
    } finally {
        saving.value = false;
    }
}

function resetForm(): void {
    form.value = {
        assetId: '',
        type: InvestmentTransactionType.Buy,
        tradeDate: new Date().toISOString().split('T')[0],
        quantity: 0,
        price: 0,
        amount: 0,
        fee: 0,
        comment: '',
    };
}

watch(dialogVisible, (val) => {
    if (val) {
        resetForm();
    }
});
</script>

<style scoped>
.asset-code {
    color: rgb(var(--v-theme-primary));
    font-weight: 500;
    margin-right: 8px;
}
</style>

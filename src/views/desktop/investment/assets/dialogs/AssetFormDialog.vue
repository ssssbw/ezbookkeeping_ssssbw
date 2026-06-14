<template>
    <v-dialog :model-value="modelValue" max-width="500" @update:model-value="$emit('update:modelValue', $event)">
        <snack-bar ref="snackbar" />
        <v-card>
            <v-toolbar color="primary">
                <v-toolbar-title>{{ mode === 'create' ? tt('Add') : tt('Edit') }}</v-toolbar-title>
            </v-toolbar>
            <v-card-text class="pa-4">
                <v-text-field v-model="form.code" :label="tt('asset.AssetCode')" density="compact" variant="outlined" :rules="[v => !!v || 'Required']" />
                <v-text-field v-model="form.name" :label="tt('asset.AssetName')" density="compact" variant="outlined" :rules="[v => !!v || 'Required']" class="mt-3" />
                <div class="d-flex ga-3 mt-3">
                    <v-select v-model="form.market" :label="tt('Market')" density="compact" variant="outlined" :items="marketOptions" style="flex: 1" />
                    <v-select v-model="form.category" :label="tt('Category')" density="compact" variant="outlined" :items="categoryOptions" style="flex: 1" />
                </div>
                <div class="d-flex ga-3 mt-3">
                    <v-text-field v-model="form.currency" :label="tt('Currency')" density="compact" variant="outlined" style="flex: 1" />
                    <v-select v-model="form.industry" :label="tt('Industry')" density="compact" variant="outlined" :items="industryFormOptions" clearable style="flex: 1" />
                </div>
                <v-text-field v-model="form.subCategory" :label="tt('asset.SubCategory')" density="compact" variant="outlined" clearable class="mt-3" />
            </v-card-text>
            <v-divider />
            <v-card-actions>
                <v-spacer />
                <v-btn variant="text" @click="$emit('update:modelValue', false)">{{ tt('Cancel') }}</v-btn>
                <v-btn color="primary" variant="tonal" :loading="saving" :disabled="!isFormValid" @click="onSubmit">{{ tt('Save') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { AssetCategory, InvestmentMarket, type AssetInfoResponse, type AssetModifyRequest } from '@/models/investment.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';


import { formatIndustry } from '../assetUtils.ts';

const props = defineProps<{
    modelValue: boolean;
    mode: 'create' | 'edit';
    asset?: AssetInfoResponse | null;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void;
    (e: 'saved'): void;
}>();

const { tt } = useI18n();

const saving = ref<boolean>(false);
type SnackBarType = InstanceType<typeof SnackBar>;
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const form = ref({
    id: '',
    code: '',
    name: '',
    market: InvestmentMarket.CN as number,
    category: AssetCategory.Equity as string,
    currency: 'CNY',
    industry: '',
    subCategory: ''
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
    title: formatIndustry(v, tt),
    value: v
}));

const isFormValid = computed(() => !!form.value.code && !!form.value.name);

watch(() => props.modelValue, (val) => {
    if (val && props.mode === 'edit' && props.asset) {
        form.value = {
            id: props.asset.id,
            code: props.asset.code,
            name: props.asset.name,
            market: props.asset.market,
            category: props.asset.category,
            currency: props.asset.currency,
            industry: props.asset.industry || '',
            subCategory: props.asset.subCategory || ''
        };
    } else if (val && props.mode === 'create') {
        form.value = {
            id: '',
            code: '',
            name: '',
            market: InvestmentMarket.CN,
            category: AssetCategory.Equity,
            currency: 'CNY',
            industry: '',
            subCategory: ''
        };
    }
});

async function onSubmit(): Promise<void> {
    if (!isFormValid.value) return;
    saving.value = true;
    try {
        if (props.mode === 'create') {
            await services.addGlobalAsset({
                code: form.value.code,
                market: form.value.market,
                name: form.value.name,
                category: form.value.category,
                currency: form.value.currency,
                industry: form.value.industry || undefined,
                subCategory: form.value.subCategory || undefined
            });
        } else {
            const req: AssetModifyRequest = {
                id: form.value.id,
                code: form.value.code,
                market: form.value.market,
                name: form.value.name,
                category: form.value.category,
                currency: form.value.currency,
                industry: form.value.industry || undefined,
                subCategory: form.value.subCategory || undefined
            };
            await services.modifyGlobalAsset(req);
        }
        emit('update:modelValue', false);
        emit('saved');
    } catch (e: unknown) {
        logger.error('Failed to save asset', e);
        const axiosErr = e as { response?: { data?: { errorMessage?: string } } };
        const msg = axiosErr.response?.data?.errorMessage;
        snackbar.value?.showError(msg || 'Failed to save asset');
    } finally {
        saving.value = false;
    }
}
</script>

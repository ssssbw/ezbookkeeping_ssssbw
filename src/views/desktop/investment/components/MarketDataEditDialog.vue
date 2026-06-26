<template>
    <v-dialog :model-value="modelValue" max-width="420" @update:model-value="$emit('update:modelValue', $event)">
        <v-card>
            <v-card-title class="d-flex align-center">
                <span class="text-h6">{{ isEdit ? tt('Edit Price') : tt('Add Price') }}</span>
                <v-spacer />
                <v-btn variant="text" size="small" :icon="mdiClose" @click="$emit('update:modelValue', false)" />
            </v-card-title>
            <v-card-text>
                <v-form ref="formRef">
                    <v-text-field
                        v-model="form.date"
                        :label="tt('Date')"
                        type="date"
                        variant="outlined"
                        density="compact"
                        class="mb-2"
                        :rules="[v => !!v || tt('This field is required')]"
                    />
                    <v-text-field
                        v-model.number="form.price"
                        :label="tt('Price')"
                        type="number"
                        variant="outlined"
                        density="compact"
                        class="mb-2"
                        :suffix="currency"
                        :rules="[v => v > 0 || tt('This field is required')]"
                    />
                    <v-text-field
                        v-model.number="form.volume"
                        :label="tt('Volume')"
                        type="number"
                        variant="outlined"
                        density="compact"
                        :suffix="tt('Volume')"
                    />
                </v-form>
            </v-card-text>
            <v-card-actions>
                <v-spacer />
                <v-btn variant="text" @click="$emit('update:modelValue', false)">{{ tt('Cancel') }}</v-btn>
                <v-btn color="primary" :loading="saving" @click="onSave">{{ tt('Save') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import { mdiClose } from '@mdi/js';

import { useI18n } from '@/locales/helpers.ts';
import services from '@/lib/services.ts';

const DIVISOR = 10000;

const props = defineProps<{
    modelValue: boolean;
    assetId: string;
    currency?: string;
    editDate?: string;
    editPrice?: number;
    editVolume?: number;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void;
    (e: 'saved'): void;
}>();

const { tt } = useI18n();

const formRef = ref();
const saving = ref(false);

const isEdit = ref(false);

const form = reactive({
    date: '',
    price: 0,
    volume: undefined as number | undefined,
});

watch(() => props.modelValue, (val) => {
    if (val) {
        if (props.editDate && props.editPrice !== undefined) {
            isEdit.value = true;
            form.date = props.editDate;
            form.price = props.editPrice / DIVISOR;
            form.volume = props.editVolume ? props.editVolume / DIVISOR : undefined;
        } else {
            isEdit.value = false;
            const now = new Date();
            form.date = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
            form.price = 0;
            form.volume = undefined;
        }
    }
});

async function onSave(): Promise<void> {
    if (!formRef.value) return;
    const { valid } = await formRef.value.validate();
    if (!valid) return;

    saving.value = true;
    try {
        const dateParts = form.date.split('-');
        const year = parseInt(dateParts[0] || '0');
        const month = parseInt(dateParts[1] || '1') - 1;
        const day = parseInt(dateParts[2] || '1');
        const dateTimestamp = Math.floor(new Date(year, month, day).getTime() / 1000);

        const priceValue = Math.round(form.price * DIVISOR);
        const volumeValue = form.volume !== undefined ? Math.round(form.volume * DIVISOR) : undefined;

        if (isEdit.value) {
            await services.modifyMarketData({
                assetId: props.assetId,
                date: dateTimestamp,
                price: priceValue,
                volume: volumeValue,
            });
        } else {
            await services.addMarketData({
                assetId: props.assetId,
                date: dateTimestamp,
                price: priceValue,
                volume: volumeValue,
            });
        }

        emit('saved');
        emit('update:modelValue', false);
    } catch {
        // error handled by service
    } finally {
        saving.value = false;
    }
}
</script>

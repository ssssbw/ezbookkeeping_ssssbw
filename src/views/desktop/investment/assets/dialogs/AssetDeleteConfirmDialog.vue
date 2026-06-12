<template>
    <v-dialog :model-value="modelValue" max-width="400" @update:model-value="$emit('update:modelValue', $event)">
        <v-card>
            <v-card-title>{{ tt('Confirm') }}</v-card-title>
            <v-card-text>
                <span class="text-body-2">{{ tt('asset.ConfirmDeleteAsset') }} "{{ asset?.name }}" ?</span>
            </v-card-text>
            <v-card-actions>
                <v-spacer />
                <v-btn variant="text" @click="$emit('update:modelValue', false)">{{ tt('Cancel') }}</v-btn>
                <v-btn color="error" variant="tonal" :loading="saving" @click="onConfirm">{{ tt('Delete') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import type { AssetInfoResponse } from '@/models/investment.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

const props = defineProps<{
    modelValue: boolean;
    asset: AssetInfoResponse | null;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void;
    (e: 'deleted'): void;
}>();

const { tt } = useI18n();

const saving = ref<boolean>(false);

async function onConfirm(): Promise<void> {
    if (!props.asset) return;
    saving.value = true;
    try {
        await services.deleteGlobalAsset({ id: props.asset.id });
        emit('update:modelValue', false);
        emit('deleted');
    } catch (e) {
        logger.error('Failed to delete asset', e);
    } finally {
        saving.value = false;
    }
}
</script>

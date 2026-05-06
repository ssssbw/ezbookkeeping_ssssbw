<template>
    <v-card :class="{ 'disabled': disabled }">
        <v-card-text class="d-flex align-center">
            <v-avatar color="grey" size="38">
                <v-icon size="24" :icon="icon" />
            </v-avatar>
            <span class="font-weight-bold ms-3">{{ title }}</span>
            <v-spacer/>
            <v-btn density="comfortable" color="default" variant="text" class="ms-2" :icon="true">
                <v-icon :icon="mdiDotsVertical" />
                <v-menu activator="parent">
                    <v-list>
                        <slot name="menus"></slot>
                    </v-list>
                </v-menu>
            </v-btn>
        </v-card-text>
        <v-card-text class="mt-1 pb-1">
            <div class="text-truncate text-h5 mb-7" :class="returnColorClass" v-if="!loading || returnAmount">{{ returnAmount }}</div>
            <v-skeleton-loader class="skeleton-no-margin mt-4 mb-8" type="text" width="120px" :loading="true" v-else-if="loading && !returnAmount"></v-skeleton-loader>
            <div class="text-truncate text-h5 mt-2 mb-7" style="padding-bottom: 2px" v-if="!loading && !returnAmount">{{ tt('No data') }}</div>
        </v-card-text>
        <v-card-text class="mt-6">
            <span class="text-caption">{{ datetime }}</span>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';

import {
    mdiDotsVertical
} from '@mdi/js';

import { computed } from 'vue';

const props = defineProps<{
    loading: boolean;
    disabled: boolean;
    icon: string;
    title: string;
    returnAmount: string;
    datetime: string;
}>();

const { tt } = useI18n();

const returnColorClass = computed(() => {
    if (!props.returnAmount) {
        return '';
    }
    if (props.returnAmount.startsWith('+')) {
        return 'text-success';
    } else if (props.returnAmount.startsWith('-')) {
        return 'text-error';
    }
    return '';
});
</script>

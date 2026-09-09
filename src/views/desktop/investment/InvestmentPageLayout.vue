<template>
    <main-page-layout investment-mode no-navbar>
        <template #content>
            <router-view :key="currentRoutePath" />
        </template>
    </main-page-layout>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();

const currentRoutePath = computed<string>(() => route.path);
</script>

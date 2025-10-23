import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { AuthResponse } from '@/models/auth_response.ts';

import { getOAuth2Provider, getOIDCCustomDisplayNames, getLoginPageTips } from '@/lib/server_settings.ts';
import { getClientDisplayVersion } from '@/lib/version.ts';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

export function useLoginPageBase(platform: 'mobile' | 'desktop') {
    const { getServerMultiLanguageConfigContent, getLocalizedOAuth2LoginText, setLanguage } = useI18n();

    const rootStore = useRootStore();
    const settingsStore = useSettingsStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const version = `${getClientDisplayVersion()}`;

    const username = ref<string>('');
    const password = ref<string>('');
    const passcode = ref<string>('');
    const backupCode = ref<string>('');
    const tempToken = ref<string>('');
    const twoFAVerifyType = ref<string>('passcode');
    const oauth2ClientSessionId = ref<string>('');

    const loggingInByPassword = ref<boolean>(false);
    const loggingInByOAuth2 = ref<boolean>(false);
    const verifying = ref<boolean>(false);

    const inputIsEmpty = computed<boolean>(() => !username.value || !password.value);
    const twoFAInputIsEmpty = computed<boolean>(() => {
        if (twoFAVerifyType.value === 'backupcode') {
            return !backupCode.value;
        } else {
            return !passcode.value;
        }
    });

    const oauth2LoginUrl = computed<string>(() => rootStore.generateOAuth2LoginUrl(platform, oauth2ClientSessionId.value));
    const oauth2LoginDisplayName = computed<string>(() => getLocalizedOAuth2LoginText(getOAuth2Provider(), getOIDCCustomDisplayNames()));
    const tips = computed<string>(() => getServerMultiLanguageConfigContent(getLoginPageTips()));

    function doAfterLogin(authResponse: AuthResponse): void {
        if (authResponse.user) {
            const localeDefaultSettings = setLanguage(authResponse.user.language);
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

            setExpenseAndIncomeAmountColor(authResponse.user.expenseAmountColor, authResponse.user.incomeAmountColor);
        }

        if (settingsStore.appSettings.autoUpdateExchangeRatesData) {
            exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
        }

        if (authResponse.notificationContent) {
            rootStore.setNotificationContent(authResponse.notificationContent);
        }
    }

    return {
        // constants
        version,
        // states
        username,
        password,
        passcode,
        backupCode,
        tempToken,
        twoFAVerifyType,
        oauth2ClientSessionId,
        loggingInByPassword,
        loggingInByOAuth2,
        verifying,
        // computed states
        inputIsEmpty,
        twoFAInputIsEmpty,
        oauth2LoginUrl,
        oauth2LoginDisplayName,
        tips,
        // functions
        doAfterLogin
    }
}

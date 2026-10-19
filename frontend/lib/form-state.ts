// Сериализуемый результат Server Action; пароль и JWT в нём не возвращаются.
export interface FormState {
  error?: string;
  success?: string;
  values?: Record<string, string>;
}

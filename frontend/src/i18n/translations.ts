// ─── Types ───
export type Lang = "ru" | "en";

// ─── Translation map ───
export const translations: Record<Lang, Record<string, string>> = {
  ru: {
    "auth.login": "Вход",
    "auth.welcome": "С возвращением",
    "auth.email": "Почта",
    "auth.password": "Пароль",
    "auth.signin": "Войти",
    "auth.loading": "Загрузка...",
    "auth.noAccount": "Нет аккаунта?",
    "auth.register": "Зарегистрироваться",
    "auth.loginVia": "Войти через",
    "auth.createAccount": "Создать аккаунт",
    "auth.name": "Имя",
    "auth.confirmPwd": "Подтвердите пароль",
    "auth.create": "Создать",
    "auth.hasAccount": "Уже есть аккаунт?",
    "auth.signinLink": "Войти",

    "profile.title": "Профиль",
    "profile.free": "Бесплатный",
    "profile.premium": "Премиум",
    "profile.userId": "ID пользователя",
    "profile.logout": "Выйти",

    "toast.welcomeBack": "С возвращением!",
    "toast.accountCreated": "Аккаунт создан! Войдите в систему.",
    "toast.error": "Что-то пошло не так",

    "val.invalidEmail": "Неверный email",
    "val.minChars": "Минимум 6 символов",
    "val.enterName": "Введите имя",
    "val.passwordsMatch": "Пароли не совпадают",
  },

  en: {
    "auth.login": "Login",
    "auth.welcome": "Welcome back",
    "auth.email": "Email",
    "auth.password": "Password",
    "auth.signin": "Sign In",
    "auth.loading": "Loading...",
    "auth.noAccount": "No account?",
    "auth.register": "Create one",
    "auth.loginVia": "Sign in with",
    "auth.createAccount": "Create account",
    "auth.name": "Name",
    "auth.confirmPwd": "Confirm password",
    "auth.create": "Create",
    "auth.hasAccount": "Already have an account?",
    "auth.signinLink": "Sign in",

    "profile.title": "Profile",
    "profile.free": "Free",
    "profile.premium": "Premium",
    "profile.userId": "User ID",
    "profile.logout": "Log out",

    "toast.welcomeBack": "Welcome back!",
    "toast.accountCreated": "Account created! Sign in to continue.",
    "toast.error": "Something went wrong",

    "val.invalidEmail": "Invalid email",
    "val.minChars": "At least 6 characters",
    "val.enterName": "Enter your name",
    "val.passwordsMatch": "Passwords don't match",
  },
};

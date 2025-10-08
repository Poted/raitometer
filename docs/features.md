

### **Faza I: MVP (Minimum Viable Product)**

**Cel główny:** Szybkie dostarczenie realnej wartości dla pierwszej grupy klientów (firm remontowych) i udowodnienie, że model biznesowy działa. Produkt ma być prosty, ale w pełni funkcjonalny i niezawodny w swoim głównym zadaniu.

**Kluczowe Funkcjonalności:**

**1. Panel Administracyjny dla Firmy (`Raite.at/dashboard`)**
    * ✅ **Rejestracja i Logowanie:** Bezpieczny system kont dla firm.
    * ✅ **Zarządzanie Usługami:** Możliwość stworzenia i edycji usług (np. "Malowanie", "Układanie Paneli").
    * ✅ **Kreator Reguł Cenowych:** Prosty interfejs, w którym firma może dla każdej usługi zdefiniować:
        * **Zmienne:** (np. metraż, stan ścian, liczba okien).
        * **Reguły:** (np. cena za metr, mnożnik za trudne warunki, stała cena za sztukę).
    * ✅ **Podgląd Unikalnego Linku:** Widoczne, łatwe do skopiowania łącze do publicznego kalkulatora firmy.

**2. Publiczny Kalkulator (Dostępny przez unikalny link)**
    * ✅ **Wizualny i Interaktywny Formularz:** Prosty, krok-po-kroku formularz dla klienta końcowego, oparty o zmienne zdefiniowane przez firmę.
    * ✅ **Generowanie Wstępnej Wyceny:** Automatyczne obliczenie szacunkowej ceny na podstawie wprowadzonych danych i reguł.
    * ✅ **Formularz Kontaktowy:** Pola na imię, e-mail i telefon.
    * ✅ **Wymagane Zgody:** Checkbox z akceptacją regulaminu i polityki prywatności (kluczowe dla RODO).

**3. System Powiadomień i Zapisu Danych**
    * ✅ **Automatyczne Maile:** Po wypełnieniu formularza system wysyła dwa maile:
        * Do klienta końcowego z podsumowaniem i informacją, że firma się skontaktuje.
        * Do firmy budowlanej z kompletem danych i szacunkową wyceną.
    * ✅ **Baza Danych Wycen:** Wszystkie zapytania są zapisywane w panelu firmy z podstawowym statusem "Nowe Zapytanie".

---

### **Faza II: "Inteligentny Asystent"**

**Cel główny:** Wprowadzenie automatyzacji opartej na AI i stworzenie kompletnego obiegu pracy, od wstępnej wyceny po finalną ofertę. To jest etap, na którym `Raite` staje się "niezbędnym" narzędziem.

**Kluczowe Funkcjonalności (dodawane do istniejącej bazy z MVP):**

**1. Analiza Zdjęć przez AI**
    * ✅ **Przesyłanie Zdjęć:** Klient końcowy może dodać zdjęcia do formularza wyceny.
    * ✅ **Integracja z API Vision AI:** Wykorzystanie gotowych modeli (np. Google Vision) do:
        * **Automatycznego Liczenia Obiektów:** System identyfikuje i liczy okna, drzwi.
        * **Wstępnego Szacowania Metrażu:** Na podstawie rozpoznanych obiektów i proporcji, system sugeruje przybliżoną powierzchnię.
    * ✅ **Wstępne Wypełnianie Formularza:** Dane z analizy AI automatycznie wypełniają pola w kalkulatorze (np. "Liczba okien: 2"). Klient może je skorygować.

**2. Mobilny Formularz Korekty dla Fachowca**
    * ✅ **Tryb "Korekty na Miejscu":** W panelu firmy, każdą wstępną wycenę można otworzyć w trybie edycji na telefonie/tablecie.
    * ✅ **Edycja Wyceny "na Żywo":** Fachowiec na miejscu koryguje dane (zmienione metraże, dodane usługi, inny stan ścian).
    * ✅ **Przeliczanie Ceny w Czasie Rzeczywistym:** Aplikacja na bieżąco aktualizuje finalną cenę na podstawie wprowadzonych korekt.

**3. Generowanie i Wysyłka Finalnej Oferty**
    * ✅ **Przycisk "Generuj i Wyślij Ofertę":** Po zakończeniu korekty, jednym kliknięciem system wykonuje akcję.
    * ✅ **Automatyczne Tworzenie PDF:** Generowanie profesjonalnie wyglądającego dokumentu PDF z finalną ofertą, zawierającego logo i dane firmy budowlanej.
    * ✅ **Wysyłka Oferty na Maile:** PDF jest automatycznie wysyłany do klienta końcowego oraz do firmy.
    * ✅ **Zaawansowane Statusy Wycen:** W panelu pojawiają się nowe statusy: "Wysłano Ofertę Finalną", "Zaakceptowana", "Odrzucona", co pozwala na lepsze zarządzanie sprzedażą.
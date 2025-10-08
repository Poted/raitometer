import os
import sys
import json
import google.generativeai as genai
from PIL import Image

# --- Konfiguracja ---
# Pobierz klucz API ze zmiennej środowiskowej.
API_KEY = os.getenv('GOOGLE_API_KEY')
if not API_KEY:
    raise ValueError("Nie znaleziono klucza API. Ustaw zmienną środowiskową GOOGLE_API_KEY.")

genai.configure(api_key=API_KEY)

# --- Krok 1: Wylistuj dostępne modele ---
# Ten kod programowo sprawdzi, jakich modeli możesz używać.
def list_available_models():
    """Wyświetla listę modeli, które wspierają metodę 'generateContent'."""
    print("🔎 Sprawdzanie dostępnych modeli dla Twojego klucza API...")
    print("-" * 30)
    found_models = False
    try:
        for m in genai.list_models():
            if 'generateContent' in m.supported_generation_methods:
                print(m.name)
                found_models = True
        if not found_models:
            print("Nie znaleziono żadnych kompatybilnych modeli.")
        else:
            print("\n✅ Powyżej znajduje się lista modeli, których możesz użyć.")
            print("Skopiuj nazwę jednego z nich i wklej ją do zmiennej MODEL_NAME poniżej.")
    except Exception as e:
        print(f"❌ Wystąpił błąd podczas pobierania listy modeli: {e}")
    print("-" * 30)


# --- Krok 2: Ustaw poprawny model i uruchom analizę ---
# WAŻNE: Upewnij się, że używasz nazwy modelu, która zadziałała dla Ciebie ostatnio.
MODEL_NAME = "models/gemini-2.0-flash-lite"
model = genai.GenerativeModel(MODEL_NAME)

# --- Główna funkcja ---
def analyze_image(image_path: str):
    """
    Analizuje podane zdjęcie i próbuje sparsować odpowiedź jako JSON.
    """
    print(f"\n📄 Analizowanie obrazu: {image_path}")
    print(f"🤖 Używany model: {MODEL_NAME}")

    if not os.path.exists(image_path):
        print(f"❌ BŁĄD: Plik obrazu nie został znaleziony w ścieżce: {image_path}")
        return

    try:
        img = Image.open(image_path)
        
        # ZMIANA: Nowy, bardziej precyzyjny prompt, który prosi o format JSON.
        prompt = [
            "Przeanalizuj to zdjęcie pod kątem elementów istotnych dla wyceny remontu.",
            "Twoim zadaniem jest zwrócić odpowiedź WYŁĄCZNIE w formacie JSON.",
            "Nie dodawaj żadnych wyjaśnień ani formatowania markdown.",
            "Stwórz listę (array) stringów, gdzie każdy string to jeden zidentyfikowany element.",
            "Przykład oczekiwanej odpowiedzi: [\"parkiet w jodełkę\", \"białe ściany\", \"okno PCV\", \"grzejnik żeberkowy\"]",
            img,
        ]

        print("🤖 Wysyłanie zapytania do Google AI (z prośbą o JSON)... Proszę czekać.")
        response = model.generate_content(prompt)

        print("\n✅ Odpowiedź z AI otrzymana pomyślnie!")
        print("-" * 30)
        print("Surowa odpowiedź modelu:")
        print(response.text)
        print("-" * 30)

        # Próba sparsowania odpowiedzi jako JSON
        try:
            # Czasami model może opakować JSON w bloki ```json ... ```
            clean_text = response.text.strip().replace("```json", "").replace("```", "").strip()
            
            parsed_json = json.loads(clean_text)
            print("🤖 Sukces! Odpowiedź to poprawny JSON.")
            print("Zidentyfikowane elementy:")
            for item in parsed_json:
                print(f"- {item}")

        except json.JSONDecodeError:
            print("🚨 Ostrzeżenie: Nie udało się sparsować odpowiedzi jako JSON.")
            print("Model nie zastosował się w pełni do instrukcji formatowania.")

    except Exception as e:
        print(f"❌ Wystąpił nieoczekiwany błąd podczas analizy obrazu: {e}")

# --- Uruchomienie skryptu ---
if __name__ == "__main__":
    if len(sys.argv) < 2:
        list_available_models()
        sys.exit(0)
    
    file_path = sys.argv[1]
    analyze_image(file_path)


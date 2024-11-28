package ui_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/security"
	"github.com/chromedp/chromedp"
)

func getBaseURL() string {
	if url := os.Getenv("APP_URL"); url != "" {
		return url
	}
	return "http://localhost:8080"
}

func takeScreenshot(_ context.Context, name string, testName string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// Crear directorio si no existe
		screenshotPath := filepath.Join("screenshots", testName)
		if err := os.MkdirAll(screenshotPath, 0755); err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}

		// Tomar screenshot usando el protocolo CDP directamente
		var buf []byte
		if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
			// Obtener el layout de la página
			var err error
			buf, err = page.CaptureScreenshot().
				WithQuality(90).
				WithClip(&page.Viewport{
					X:      0,
					Y:      0,
					Width:  1280, // Ancho fijo
					Height: 900,  // Alto fijo
					Scale:  1,
				}).Do(ctx)
			return err
		})); err != nil {
			return err
		}

		filename := filepath.Join(screenshotPath, fmt.Sprintf("%s.png", name))
		if err := os.WriteFile(filename, buf, 0644); err != nil {
			return fmt.Errorf("error writing screenshot: %v", err)
		}

		return nil
	})
}

func setupChrome(_ *testing.T) (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/chromium"),
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("ignore-certificate-errors", true),
		// Agregar flags adicionales para screenshots
		chromedp.WindowSize(1280, 900),
		chromedp.Flag("force-device-scale-factor", "1"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel2 := chromedp.NewContext(allocCtx)

	cancelFunc := func() {
		cancel2()
		cancel()
	}

	// Configurar timeouts más largos
	ctx, cancel3 := context.WithTimeout(ctx, 60*time.Second)
	go func() {
		<-ctx.Done()
		cancel3()
	}()

	chromedp.Run(ctx, security.Enable())
	chromedp.Run(ctx, security.SetIgnoreCertificateErrors(true))

	return ctx, cancelFunc
}

func TestConnectionToApp(t *testing.T) {
	start := time.Now()
	baseURL := getBaseURL()
	t.Logf("TEST: Conexión inicial a %s", baseURL)

	ctx, cancel := setupChrome(t)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(baseURL),
		chromedp.Sleep(2*time.Second),
		takeScreenshot(ctx, "initial-connection", "TestConnectionToApp"),
	)

	if err != nil {
		t.Fatalf("❌ Error de conexión: %v", err)
	}
	t.Logf("✅ Conexión exitosa (%.2fs)", time.Since(start).Seconds())
}

func TestPostsUI(t *testing.T) {
	baseURL := getBaseURL()
	t.Logf("TEST SUITE: UI Posts en %s", baseURL)

	ctx, cancel := setupChrome(t)
	defer cancel()

	t.Run("Carga de página", func(t *testing.T) {
		start := time.Now()
		var title, h1Text string
		var attempts int

		for attempts = 1; attempts <= 5; attempts++ {
			err := chromedp.Run(ctx,
				chromedp.Navigate(baseURL+"/web/posts"),
				chromedp.Sleep(2*time.Second),
				chromedp.Title(&title),
				chromedp.Text("h1", &h1Text),
				takeScreenshot(ctx, fmt.Sprintf("page-load-attempt-%d", attempts), "PageLoad"),
			)

			if err == nil {
				break
			}
			time.Sleep(2 * time.Second)
		}

		if title != "Super Descuentos" {
			t.Errorf("❌ Título incorrecto: esperado 'Super Descuentos', obtenido '%s'", title)
		}
		t.Logf("✅ Página cargada en intento %d (%.2fs)", attempts, time.Since(start).Seconds())
	})

	t.Run("Crear post", func(t *testing.T) {
		start := time.Now()
		tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04")

		err := chromedp.Run(ctx,
			chromedp.Navigate(baseURL+"/web/posts"),
			chromedp.WaitVisible("#create-post-form", chromedp.ByID),
			takeScreenshot(ctx, "before-create", "CreatePost"),
			chromedp.SetValue("#title", "Test Post from UI", chromedp.ByID),
			chromedp.SetValue("#description", "This is a test post created by UI tests", chromedp.ByID),
			chromedp.SetValue("#url", "https://example.com", chromedp.ByID),
			chromedp.SetValue("#expire-time", tomorrow, chromedp.ByID),
			takeScreenshot(ctx, "filled-form", "CreatePost"),
			chromedp.Click("button[type='submit']"),
			chromedp.Sleep(2*time.Second),
			takeScreenshot(ctx, "after-create", "CreatePost"),
		)

		if err != nil {
			t.Fatalf("❌ Error al crear post: %v", err)
		}
		t.Logf("✅ Post creado: 'Test Post from UI' (%.2fs)", time.Since(start).Seconds())
	})

	t.Run("Eliminar post", func(t *testing.T) {
		start := time.Now()
		err := chromedp.Run(ctx,
			chromedp.Navigate(baseURL+"/web/posts"),
			chromedp.WaitVisible(".post"),
			takeScreenshot(ctx, "before-delete", "DeletePost"),
			chromedp.Evaluate(`window.confirm = function() { return true; }`, nil),
			chromedp.Click("button:contains('Eliminar')"),
			chromedp.Sleep(2*time.Second),
			takeScreenshot(ctx, "after-delete", "DeletePost"),
		)

		if err != nil {
			t.Fatalf("❌ Error al eliminar post: %v", err)
		}
		t.Logf("✅ Post eliminado (%.2fs)", time.Since(start).Seconds())
	})

	t.Run("Validación de formulario", func(t *testing.T) {
		start := time.Now()
		err := chromedp.Run(ctx,
			chromedp.Navigate(baseURL+"/web/posts"),
			chromedp.WaitVisible("#create-post-form", chromedp.ByID),
			takeScreenshot(ctx, "before-validation", "FormValidation"),
			chromedp.Click("button[type='submit']"),
			chromedp.WaitVisible("#title:invalid", chromedp.ByQuery),
			takeScreenshot(ctx, "after-validation", "FormValidation"),
		)

		if err != nil {
			t.Fatalf("❌ Error en validación del formulario: %v", err)
		}
		t.Logf("✅ Validación de formulario correcta (%.2fs)", time.Since(start).Seconds())
	})
}

import { test, expect } from '@playwright/test';

test.use({
  ignoreHTTPSErrors: true,
});

// Configuration - adjust these to match your application's setup
const BASE_URL = process.env.BASE_URL || 'http://localhost:8080';

test.describe('Super Descuentos Application E2E Tests', () => {
  // Test creating a new author
  test('create a new author', async ({ page }) => {
    await page.goto(`${BASE_URL}/authors`);
    await page.screenshot({path: 'screenshots/authors.png'});

    // Fill out author creation form
    await page.fill('input[name="name"]', 'Test Author');
    await page.fill('input[name="email"]', 'testauthor@example.com');
    
    await page.click('button[type="submit"]');

    // Wait for navigation or confirmation
    await page.waitForURL(`${BASE_URL}/authors`);

    // Verify author was created (you might need to adjust this selector)
    const authorExists = await page.getByText('testauthor@example.com').isVisible();
    expect(authorExists).toBeTruthy();
  });

  // Test creating a new post
  test('create a new post', async ({ page }) => {
    await page.goto(`${BASE_URL}/posts`);
    await page.screenshot({path: 'screenshots/posts.png'});

    // First, ensure we have an author to select
    await page.goto(`${BASE_URL}/authors`);
    await page.fill('input[name="name"]', 'Post Author');
    await page.fill('input[name="email"]', 'postauthor@example.com');
    await page.click('button[type="submit"]');

    // Navigate to posts page
    await page.goto(`${BASE_URL}/posts`);

    // Fill out post creation form
    await page.selectOption('select[name="author_email"]', 'postauthor@example.com');
    await page.fill('input[name="title"]', 'Test Discount Offer');
    await page.fill('textarea[name="content"]', 'Amazing discount available now!');
    await page.fill('input[name="url"]', 'https://example.com/discount');
    
    await page.click('button[type="submit"]');

    // Wait for navigation
    await page.waitForURL(`${BASE_URL}/`);

    // Verify post was created and appears on home page
    const postTitle = await page.getByText('Test Discount Offer').first();
    expect(postTitle).toBeVisible();
  });

  // Test home page functionality
  test('home page displays posts', async ({ page }) => {
    // Ensure at least one post exists
    await page.goto(`${BASE_URL}/posts`);
    await page.screenshot({path: 'screenshots/posts.png'});
    
    // If no posts exist, create one
    const noPostsMessage = await page.getByText('No se encontraron posts en el sistema').isVisible();
    if (noPostsMessage) {
      // Create an author first
      await page.goto(`${BASE_URL}/authors`);
      await page.fill('input[name="name"]', 'Home Page Author');
      await page.fill('input[name="email"]', 'homeauthor@example.com');
      await page.click('button[type="submit"]');

      // Create a post
      await page.goto(`${BASE_URL}/posts`);
      await page.selectOption('select[name="author_email"]', 'homeauthor@example.com');
      await page.fill('input[name="title"]', 'Home Page Test Post');
      await page.fill('textarea[name="content"]', 'Test content for home page');
      await page.fill('input[name="url"]', 'https://example.com/test');
      await page.click('button[type="submit"]');
    }

    // Navigate to home page and verify posts
    await page.goto(`${BASE_URL}/`);
    const postCards = await page.locator('.card').count();
    expect(postCards).toBeGreaterThan(0);
  });

  // Test external link functionality
  test('post external link works', async ({ page, context }) => {
    await page.goto(`${BASE_URL}/`);
    await page.screenshot({path: 'screenshots/home.png'});

    // Find the first external link and open it in a new page
    const firstExternalLink = await page.locator('.card-link').first();
    const newPagePromise = context.waitForEvent('page');
    await firstExternalLink.click();
    const newPage = await newPagePromise;

    // Wait for the new page to load
    await newPage.waitForLoadState('load');

    // Verify the new page is not empty
    const url = newPage.url();
    expect(url).not.toBe('about:blank');
  });
});
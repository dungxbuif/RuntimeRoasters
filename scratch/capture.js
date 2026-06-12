const { chromium } = require('playwright');
const fs = require('fs');

(async () => {
  // Ensure assets directory exists
  if (!fs.existsSync('docs/assets')) {
    fs.mkdirSync('docs/assets', { recursive: true });
  }

  const browser = await chromium.launch();
  const context = await browser.newContext({
    recordVideo: {
      dir: 'docs/assets', // Video output directory
      size: { width: 1280, height: 720 }
    },
    viewport: { width: 1280, height: 720 }
  });

  const page = await context.newPage();

  console.log("Navigating to localhost:3000...");
  try {
    await page.goto('http://localhost:3000', { waitUntil: 'networkidle', timeout: 30000 });
  } catch (e) {
    console.error("Failed to load page. Is the server running?", e);
    await browser.close();
    process.exit(1);
  }

  // Wait a few seconds for data to load and render
  console.log("Waiting for rendering...");
  await page.waitForTimeout(3000);

  console.log("Taking screenshot...");
  await page.screenshot({ path: 'docs/assets/demo-screenshot.png' });

  // Wait a little more to record some video action if any
  await page.waitForTimeout(2000);

  // Close context to finish recording video
  await context.close();

  // Rename the video file from a random hash to a readable name
  const files = fs.readdirSync('docs/assets');
  const videoFile = files.find(f => f.endsWith('.webm'));
  if (videoFile) {
    fs.renameSync(`docs/assets/${videoFile}`, 'docs/assets/demo-video.webm');
    console.log("Video saved to docs/assets/demo-video.webm");
  }

  await browser.close();
  console.log("Done!");
})();

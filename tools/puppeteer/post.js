#!/usr/bin/env node
/**
 * TormentNexus — Semi-Automated Poster
 * Opens browser, fills in content, waits for user review before posting
 * 
 * Usage:
 *   node post-to-reddit.js <subreddit>
 *   node post-to-hn.js
 *   node post-to-twitter.js
 */

const puppeteer = require("puppeteer");
const fs = require("fs");
const path = require("path");

const AUTO_DIR = path.join(__dirname, "..", "..", "marketing", "auto");

// Find the latest JSON file matching pattern
function findLatestFile(pattern) {
  const files = fs.readdirSync(AUTO_DIR)
    .filter(f => f.endsWith(".json") && f.includes(pattern))
    .sort()
    .reverse();
  
  if (files.length === 0) {
    console.error(`No files found matching: ${pattern}`);
    process.exit(1);
  }
  
  return path.join(AUTO_DIR, files[0]);
}

// Read JSON content
function readContent(file) {
  const data = JSON.parse(fs.readFileSync(file, "utf8"));
  return data;
}

// Wait for user input
function waitForEnter(message) {
  return new Promise(resolve => {
    console.log(message);
    process.stdin.once("data", () => resolve());
  });
}

// Post to Reddit
async function postToReddit(subreddit) {
  console.log(`\n=== Posting to r/${subreddit} ===\n`);
  
  const file = findLatestFile(`reddit_${subreddit}`);
  const content = readContent(file);
  
  console.log("Title:", content.title);
  console.log("Body preview:", content.body.substring(0, 200) + "...");
  console.log("");
  
  const browser = await puppeteer.launch({
    headless: false,
    args: ["--no-sandbox", "--disable-setuid-sandbox", "--start-maximized"]
  });
  
  const page = await browser.newPage();
  await page.setViewport({ width: 1280, height: 900 });
  
  // Navigate to Reddit submit page
  console.log("Opening Reddit...");
  await page.goto(`https://www.reddit.com/r/${subreddit}/submit`, {
    waitUntil: "networkidle2",
    timeout: 30000
  });
  
  // Wait for user to log in if needed
  await waitForEnter("\n>>> Log in to Reddit if needed, then press ENTER to continue...");
  
  // Select "Text" post type
  try {
    const textButton = await page.$('button:has-text("Post"), [data-testid="text-post"]');
    if (textButton) await textButton.click();
  } catch (e) {
    console.log("Note: You may need to manually select 'Text' post type");
  }
  
  // Fill in title
  console.log("Filling in title...");
  const titleInput = await page.$('textarea[placeholder*="Title"], input[placeholder*="Title"]');
  if (titleInput) {
    await titleInput.click();
    await titleInput.type(content.title.substring(0, 300));
  }
  
  // Fill in body
  console.log("Filling in body...");
  const bodyInput = await page.$('textarea[placeholder*="body"], [contenteditable="true"]');
  if (bodyInput) {
    await bodyInput.click();
    // Type body in chunks to avoid issues
    const chunks = content.body.match(/.{1,500}/g) || [content.body];
    for (const chunk of chunks) {
      await page.keyboard.type(chunk, { delay: 10 });
    }
  }
  
  console.log("\n✅ Content filled in!");
  await waitForEnter("\n>>> Review the post in the browser. Press ENTER to submit...");
  
  // Click submit
  const submitButton = await page.$('button:has-text("Post"), button[type="submit"]');
  if (submitButton) {
    await submitButton.click();
    console.log("✅ Post submitted!");
  } else {
    console.log("⚠️ Could not find submit button. Please click it manually.");
  }
  
  await waitForEnter("\n>>> Press ENTER to close the browser...");
  await browser.close();
}

// Post to Hacker News
async function postToHN() {
  console.log("\n=== Posting to Hacker News ===\n");
  
  const file = findLatestFile("hackernews");
  const content = readContent(file);
  
  // Parse the nested JSON
  let title, body;
  try {
    const parsed = JSON.parse(content.body);
    title = parsed.title || content.title;
    body = parsed.body || content.body;
  } catch {
    title = content.title;
    body = content.body;
  }
  
  console.log("Title:", title);
  console.log("Body preview:", body.substring(0, 200) + "...");
  console.log("");
  
  const browser = await puppeteer.launch({
    headless: false,
    args: ["--no-sandbox", "--disable-setuid-sandbox", "--start-maximized"]
  });
  
  const page = await browser.newPage();
  await page.setViewport({ width: 1280, height: 900 });
  
  // Navigate to HN submit page
  console.log("Opening Hacker News...");
  await page.goto("https://news.ycombinator.com/submit", {
    waitUntil: "networkidle2",
    timeout: 30000
  });
  
  // Wait for user to log in if needed
  await waitForEnter("\n>>> Log in to HN if needed, then press ENTER to continue...");
  
  // Fill in title
  console.log("Filling in title...");
  const titleInput = await page.$('input[name="title"]');
  if (titleInput) {
    await titleInput.click();
    await titleInput.type(title.substring(0, 80));
  }
  
  // Fill in URL
  console.log("Filling in URL...");
  const urlInput = await page.$('input[name="url"]');
  if (urlInput) {
    await urlInput.click();
    await urlInput.type("https://github.com/MDMAtk/TormentNexus");
  }
  
  // Fill in text
  console.log("Filling in text...");
  const textInput = await page.$('textarea[name="text"]');
  if (textInput) {
    await textInput.click();
    const chunks = body.match(/.{1,500}/g) || [body];
    for (const chunk of chunks) {
      await page.keyboard.type(chunk, { delay: 10 });
    }
  }
  
  console.log("\n✅ Content filled in!");
  await waitForEnter("\n>>> Review the post in the browser. Press ENTER to submit...");
  
  // Click submit
  const submitButton = await page.$('input[type="submit"]');
  if (submitButton) {
    await submitButton.click();
    console.log("✅ Post submitted!");
  }
  
  await waitForEnter("\n>>> Press ENTER to close the browser...");
  await browser.close();
}

// Post to Twitter
async function postToTwitter() {
  console.log("\n=== Posting Twitter Thread ===\n");
  
  const file = findLatestFile("twitter");
  const content = readContent(file);
  
  console.log("Thread:");
  content.tweets.forEach((t, i) => console.log(`  ${i + 1}. ${t}`));
  console.log("");
  
  const browser = await puppeteer.launch({
    headless: false,
    args: ["--no-sandbox", "--disable-setuid-sandbox", "--start-maximized"]
  });
  
  const page = await browser.newPage();
  await page.setViewport({ width: 1280, height: 900 });
  
  // Navigate to Twitter
  console.log("Opening Twitter...");
  await page.goto("https://twitter.com/compose/tweet", {
    waitUntil: "networkidle2",
    timeout: 30000
  });
  
  // Wait for user to log in
  await waitForEnter("\n>>> Log in to Twitter if needed, then press ENTER to continue...");
  
  for (let i = 0; i < content.tweets.length; i++) {
    console.log(`\nPosting tweet ${i + 1}/${content.tweets.length}...`);
    console.log(`"${content.tweets[i].substring(0, 50)}..."`);
    
    // Type tweet
    const tweetInput = await page.$('[data-testid="tweetTextarea_0"], [contenteditable="true"]');
    if (tweetInput) {
      await tweetInput.click();
      await page.keyboard.type(content.tweets[i], { delay: 20 });
    }
    
    if (i < content.tweets.length - 1) {
      await waitForEnter(`>>> Review tweet ${i + 1}. Press ENTER to post and continue...`);
      
      // Click tweet button
      const tweetButton = await page.$('[data-testid="tweetButtonInline"]');
      if (tweetButton) await tweetButton.click();
      
      // Wait for tweet to post
      await new Promise(r => setTimeout(r, 3000));
      
      // Start new tweet
      const newTweetButton = await page.$('[data-testid="tweetButton"]');
      if (newTweetButton) await newTweetButton.click();
    } else {
      await waitForEnter(`>>> Review final tweet. Press ENTER to post...`);
      const tweetButton = await page.$('[data-testid="tweetButtonInline"]');
      if (tweetButton) await tweetButton.click();
    }
  }
  
  console.log("\n✅ Thread posted!");
  await waitForEnter("\n>>> Press ENTER to close the browser...");
  await browser.close();
}

// Main
const args = process.argv.slice(2);
const command = args[0];

switch (command) {
  case "reddit":
    postToReddit(args[1] || "MachineLearning");
    break;
  case "hn":
    postToHN();
    break;
  case "twitter":
    postToTwitter();
    break;
  default:
    console.log("Usage:");
    console.log("  node post.js reddit <subreddit>");
    console.log("  node post.js hn");
    console.log("  node post.js twitter");
    console.log("");
    console.log("Examples:");
    console.log("  node post.js reddit MachineLearning");
    console.log("  node post.js reddit LocalLLaMA");
    console.log("  node post.js hn");
    console.log("  node post.js twitter");
}

# insta-tools

_A command-line tool to interact with Instagram’s API, retrieve user data, followers and following lists (and much more in the future)._

I do not recomend using this tool for any kind of scraping or data collection as it probably violates Instagram's terms of service. I must state that this tool is intended for educational purposes only.
So do not blame for any reprisal you may suffer from using this tool, for this is fully your responsability.

---

## **🚀 Installation**

### **1. Install via `go install`**

To install `insta-tools`, run:

```sh
go install github.com/Rfluid/insta-tools@latest
```

This will install the CLI tool globally on your system.

### **2. Manual Installation**

Clone the repository:

```sh
git clone https://github.com/Rfluid/insta-tools.git
cd insta-tools
```

Build the project:

```sh
go build -o bin/insta-tools .
```

Now, you can run `insta-tools` from the `bin/` directory.

---

## **🔑 Required Authentication - How to Get Cookies**

To use `insta-tools`, you need Instagram session cookies. Choose the method that suits you best. I generally go with the first.

### **📌 Method 1: Using `document.cookie` in the Console**

1. Open **Instagram** ([https://www.instagram.com](https://www.instagram.com)) and log in.
2. Open **Developer Tools**:
   - **Chrome**: Press `F12` or `Ctrl + Shift + I` (`Cmd + Opt + I` on Mac).
   - **Firefox**: Press `F12` or `Ctrl + Shift + I`.
3. Click on the **Console** tab.
4. Copy and paste this command, then press **Enter**:

   ```js
   console.log(document.cookie);
   ```

5. The output will show all your cookies:

   ```sh
   datr=abc123; ig_did=xyz456; csrftoken=token123; sessionid=987654321%3Aabcdef%3A12%3Aabcxyz;
   ```

### **📌 Method 2: Open Instagram on Your Browser**

1. Log in to [Instagram](https://www.instagram.com/).
2. Open **Developer Tools**:
   - Chrome: Press `F12` or `Ctrl + Shift + I` (`Cmd + Opt + I` on Mac).
   - Firefox: Press `F12` or `Ctrl + Shift + I`.
3. Go to the **Network** tab and refresh the page.
4. Look for a request to `www.instagram.com/api/v1/...`
5. Copy the **Cookies** header. It should contain values like:

   ```sh
   datr=...; ig_did=...; csrftoken=...; sessionid=...
   ```

6. Use this string as your `--cookies` value.

---

### **📌 Method 3: Extract Only the Necessary Cookies**

If you only need **session-related cookies**, run:

```js
console.log(
  "csrftoken=" +
    document.cookie.match(/csrftoken=([^;]+)/)[1] +
    "; " +
    "sessionid=" +
    document.cookie.match(/sessionid=([^;]+)/)[1],
);
```

📌 **Output Example**:

```sh
csrftoken=token123; sessionid=987654321%3Aabcdef%3A12%3Aabcxyz;
```

This ensures you only extract **the cookies needed** for authentication.

---

### **📌 Using These Cookies in `insta-tools`**

Once you have the cookies, use them in `insta-tools`:

```sh
insta-tools followers <userID> <count> "" --cookies "csrftoken=your_token; sessionid=your_session_id"
```

This method ensures you **retrieve valid Instagram session cookies** easily. 🚀 Let me know if you need more details!

---

## **🛠️ Usage**

### **1. Retrieve user data**

#### **Basic Example**

```sh
insta-tools user <username> "<your_cookies>"
```

Example:

```sh
insta-tools user zuck --cookies "sessionid=YOUR_SESSION_ID; csrftoken=YOUR_CSRFTOKEN"
```

- `username`: Instagram username.

#### **Save user to a File**

```sh
insta-tools user zuck -o ./users/zuck.json --cookies "sessionid=YOUR_SESSION_ID; csrftoken=YOUR_CSRFTOKEN"
```

- The output will be saved to `./users/zuck.json`.

---

### **2. Retrieve Followers**

#### **Basic Example**

```sh
insta-tools followers <userID> <count> <maxID> --cookies "<your_cookies>"
```

Example:

```sh
insta-tools followers 314216 12 "" --cookies "sessionid=YOUR_SESSION_ID; csrftoken=YOUR_CSRFTOKEN"
```

- `userID`: Instagram user ID. You can get it by using the previous command. E. g. when searching for [zuck](instagram.com/zuck)'s data, we find his user ID `314216` at `data.user.id`. You can just hit

```sh
insta-tools user zuck --cookies "<your_cookies>" | jq -r '.data.user.id'
```

to directly retrieve the user ID.

- `count`: Number of followers per request.
- `maxID`: Used for pagination (set to `""` for the first request).

#### **Retrieve All Followers**

```sh
insta-tools followers <userID> <count> <maxID> --all --threads 4 --sleep 2 --cookies "<your_cookies>"
```

- `--all`: Fetch all followers, paginating automatically.
- `--threads`: Number of concurrent API requests.
- `--sleep`: Delay (in seconds) between requests to prevent rate limits.

#### **Save Followers to a File**

```sh
insta-tools followers 314216 12 "" --all -o ./test-data/followers.json --cookies "sessionid=YOUR_SESSION_ID; csrftoken=YOUR_CSRFTOKEN"
```

- The output will be saved to `./test-data/followers.json`.

---

### **3. Retrieve Following**

#### **Basic Example**

```sh
insta-tools following <userID> <count> <maxID> --cookies "<your_cookies>"
```

Example:

```sh
insta-tools following 314216 12 "" --cookies "sessionid=YOUR_SESSION_ID; csrftoken=YOUR_CSRFTOKEN"
```

#### **Retrieve All Following**

```sh
insta-tools following <userID> <count> <maxID> --all --threads 4 --sleep 2 --cookies "<your_cookies>"
```

---

## **⚙️ Global Flags**

These flags work with all commands:

| Flag           | Description                         |
| -------------- | ----------------------------------- |
| `--cookies`    | Set Instagram session cookies       |
| `--output, -o` | Save results to a file              |
| `--threads`    | Number of concurrent API requests   |
| `--logs`       | Enable logging for better debugging |

---

## **📌 Example: Retrieve & Save Followers**

```sh
insta-tools followers 314216 12 "" --all --threads 4 --sleep 2 --cookies "sessionid=YOUR_SESSION_ID; csrftoken=YOUR_CSRFTOKEN" -o followers.json
```

This:

- Fetches **all followers**.
- Uses **4 threads**.
- **Waits 2 seconds** between requests.
- **Saves results to `followers.json`**.

---

## **💡 Troubleshooting**

### **1. `Error: received status code 403`**

**Solution:** Check if:

- Your **session cookies are valid**.
- You’re **not rate-limited** (reduce `--threads` or increase `--sleep`).

---

## **📜 License**

This project is licensed under the **MIT License**.

---

## **🙌 Contributing**

Pull requests are welcome! If you find a bug or want to add features, feel free to contribute.

🚀 **Enjoy using `insta-tools`!** Let me know if you need any modifications! 🚀

---

σΔγ

# u

A collection of utility functions for [Alfred Workflow](https://www.alfredapp.com/workflows/).  
Useful for developers, power users, and anyone who wants to work faster.  
Feel free to use, tweak, and contribute!

---

## ⚙️ Functions

### 🔹 UUID Generator — `u uuid`

Generates a random UUID (Universally Unique Identifier), e.g.:
```
1b491dd5-194c-4d10-b16e-49abc6b6c882
```

- Takes no input
- Useful for generating unique IDs for items, logs, tasks, etc.
- Powered by [`google/uuid`](https://github.com/google/uuid)

---

### 🕒 Timestamp & Duration Utility — `u ts`

Easily work with **timestamps**, **durations**, and **custom formatting**.

#### ✅ Supports
- Zero input
- One input (unary)
- Two inputs (binary)
- ❌ No more than 2 inputs allowed

---

### 📥 Input Types

#### 📌 Timestamp
Accepts:
- Epoch seconds: `1625097600`
- Milliseconds: `1625097600123`
- Microseconds: `1625097600123456`
- Predefined format string: `'2025-06-30 12:34:56'` (wrap in single quotes)

Outputs:
- Milliseconds
- Configured formats (e.g. ISO8601, RFC3339 — configurable via Alfred)

---

#### 📌 Duration
Supports:
- ISO 8601: `P1Y2M3DT4H5M6S`
- Shorthand: `1y`, `3mo`, `10d`, `2h`, `45m`, `30s`
- Combine with `+` or `-` operator

Powered by: [`sosodev/duration`](https://github.com/sosodev/duration)

---

#### 📌 Format Strings
Format current or given time using:
- `strftime` (e.g. `'%Y-%m-%d %H:%M'`)
- Go time format (e.g. `'2006-01-02 15:04:05'`)

---

### 💡 Examples

```bash
u ts
# → Current time in ms + default formats

u ts 1650000000000
# → Convert milliseconds to readable timestamp

u ts 1h
# → Convert duration to readable form

u ts '2025-06-30 10:00:00' + 3h
# → 2025-06-30 13:00:00

u ts '2006-01-02','%Y-%m-%d %H:%M'
# → Format now using both styles
```

---

## ❣️ Inspirations
Icon by [Web development icons created by srip - Flaticon](https://www.flaticon.com/free-icons/web-development)

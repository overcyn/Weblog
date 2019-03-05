# Weblog

Write a log to database
`POST https://overcyn.appspot.com?message=Body&date=2006-01-02T15:04:05-07:00`

Clear all logs
`DELETE https://overcyn.appspot.com`

### Run dev server
dev_appserver.py app.yaml

### Deploy app
gcloud app deploy

### Integrating into iOS App

```
// Debug code... don't commit

fileprivate let dispatchQueue = DispatchQueue(label: "KDLog")

public func KDLog(_ message: String) {
    var body: String = "message=\(message)"
    if #available(iOS 11.0, watchOS 4.0, *) {
        let dateString = ISO8601DateFormatter.string(from: Date(), timeZone: TimeZone(abbreviation: "GMT")!, formatOptions: [.withInternetDateTime, .withFractionalSeconds])
        body = "date=\(dateString)&message=\(message)"
    }
    print(body)
    
    dispatchQueue.async {
        KDPush(body)
        KDSyncLog()
    }
}

func KDSyncLog() {
    let semaphore = DispatchSemaphore(value: 0)
    ProcessInfo.processInfo.performExpiringActivity(withReason: "KDSyncLog") { expiring in
        if expiring {
            semaphore.signal()
        } else {
            dispatchQueue.sync {
                KDSyncLogIter {
                    semaphore.signal()
                }
            }
            _ = semaphore.wait()
        }
    }
}

func KDSyncLogIter(_ completion: @escaping () -> ()) {
    if let message = KDPop() {
        var request: URLRequest = URLRequest(url: URL(string: "https://overcyn.appspot.com")!)
        request.httpMethod = "POST"
        request.httpBody = message.data(using: .utf8)
        URLSession.shared.dataTask(with: request) {_, _, err in
            dispatchQueue.async {
                if err != nil {
                    KDPush(message)
                } else {
                    KDSyncLogIter(completion)
                }
            }
        }.resume()
    } else {
        completion()
    }
}

func KDPop() -> String? {
    var messages = UserDefaults.standard.stringArray(forKey: "KDLog") ?? []
    let message = messages.popLast()
    UserDefaults.standard.set(messages, forKey: "KDLog")
    return message
}

func KDPush(_ message: String) {
    var messages = UserDefaults.standard.stringArray(forKey: "KDLog") ?? []
    messages.append(message)
    UserDefaults.standard.set(messages, forKey: "KDLog")
}
// Debug code... Don't commit
```

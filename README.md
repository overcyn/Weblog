# Weblog

> env GOOS=linux GOARCH=amd64 go build .
> cd ..
> scp -r Weblog overcyn@173.255.221.39:/home/overcyn
> ssh overcyn@173.255.221.39
> cd Weblog
> ./Weblog &

Write a log to database

`POST http://173.255.221.39:8002?message=Body&date=2006-01-02T15:04:05-07:00`

Clear all logs

`DELETE http://173.255.221.39:8002`

### Run dev server
dev_appserver.py app.yaml

### Deploy app
gcloud app deploy

### Integrating into iOS App

```
// Debug code... don't commit

fileprivate let dispatchQueue = DispatchQueue(label: "KDLog")

public func KDLog(_ message: String) {
    let dateString = ISO8601DateFormatter.string(from: Date(), timeZone: TimeZone(abbreviation: "GMT")!, formatOptions: [.withInternetDateTime, .withFractionalSeconds])
    let body = "date=\(dateString)&message=\(message)"
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
            semaphore.wait()
        }
    }
}

func KDSyncLogIter(_ completion: @escaping () -> ()) {
    if let message = KDPop() {
        var request: URLRequest = URLRequest(url: URL(string: "http://173.255.221.39:8002")!)
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



class Solution {
private:
    // data members
    unordered_set<string> seen_urls;
    // urls to crawl
    // add url to the seen set before pushing into the q
    queue<string> q;  
    string host_name;

    // locks the q, seen_urls
    mutex m;  
    condition_variable cv;

    // exit loop condition:
    // no worker working, and queue is empty
    int num_active_workers = 0;
    bool done = false;


public:
    string GetHostName(const string& url) {
        size_t start = url.find('/');
        start += 2;
        size_t end = url.find('/', start);

        return url.substr(start, end-start);
    }

    // wait until some work exist in q
    // add results into seen_urls
    void worker_thread(HtmlParser* htmlParser) {
        while (true) {
            unique_lock<mutex> ul(m);
            // wake up when notified and q is not empty
            cv.wait(ul, [&]() { return !q.empty() || done; });
            if (done) {
                break;
            }

            // now we have the lock
            ++num_active_workers;
            string url = q.front();
            q.pop();
            // values in the queue are valid.

            // crawl the new url
            ul.unlock();
            auto new_urls = htmlParser->getUrls(url);

            // add crawled urls into data members
            ul.lock();
            for (string& val : new_urls) {
                if (seen_urls.count(val)==1 || GetHostName(val) != host_name) {
                    continue;
                }
                seen_urls.insert(val);
                q.push(val);
            }
            --num_active_workers;

            if (q.empty() && num_active_workers == 0) {
                done = true;
            }

            cv.notify_all();
        }

    }

    vector<string> crawl(string startUrl, HtmlParser& htmlParser) {
        host_name = GetHostName(startUrl);
        // cout << thread::hardware_concurrency();

        vector<thread> workers;
        int num_threads = 5;  // configurable
        for (int i = 0; i < num_threads; ++i) {
            workers.emplace_back(&Solution::worker_thread, this, &htmlParser);
        }

        q.push(startUrl);
        seen_urls.insert(startUrl);
        cv.notify_one();

        for (thread& t : workers) {
            t.join();
        }

        return vector<string>(seen_urls.begin(), seen_urls.end());
    }
};

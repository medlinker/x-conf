package com.medlinker.xconf;

import java.io.BufferedOutputStream;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.net.ConnectException;
import java.net.InetAddress;
import java.net.URI;
import java.net.UnknownHostException;
import java.util.HashMap;
import java.util.Map;
import java.util.Map.Entry;
import java.util.Properties;
import java.util.concurrent.TimeoutException;

import mousio.client.retry.RetryWithTimeout;
import mousio.etcd4j.EtcdClient;
import mousio.etcd4j.promises.EtcdResponsePromise;
import mousio.etcd4j.responses.EtcdAuthenticationException;
import mousio.etcd4j.responses.EtcdException;
import mousio.etcd4j.responses.EtcdKeysResponse;
import mousio.etcd4j.responses.EtcdKeysResponse.EtcdNode;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * x-conf java client
 * 
 * @author mac
 *
 */
public class Client {
    
    private final Logger LOG = LoggerFactory.getLogger(Client.class);
    
    private Map<String, String> entries = new HashMap<>(64);

    private EtcdClient ec;

    private Properties prop = new Properties();

    private String prjName;

    private String env;

    public Client() {
        this.init();
    }

    private void init() {
        try {
            prop.load(new FileInputStream("x-conf.conf"));
            String[] urls =
                    prop.getProperty("etcd_clinet_urls", "http://127.0.0.1:2379").split(",");
            URI[] uris = new URI[urls.length];
            for (int i = 0; i < urls.length; i++) {
                uris[i] = URI.create(urls[i]);
            }
            ec = new EtcdClient(uris);
            ec.setRetryHandler(new RetryWithTimeout(200, 20000));

            prjName = prop.getProperty("prjName", "default");
            env = prop.getProperty("env", "prod");

            this.pullAll();
            this.dump();

            String info = getInfo();
            System.out.println(info);
            new Thread(() -> {
                while (true) {
                    try {
                        Thread.sleep(5000);
                        heart(info);
                    } catch (Exception e) {
                        LOG.error(e.getMessage());
                    }
                }
            }).start();
        } catch (Exception e) {
            LOG.error(e.getMessage());
            if (e instanceof TimeoutException || e instanceof ConnectException) {
                this.readFromDump();
            }
        }
    }

    private void pullAll() throws IOException, EtcdException, EtcdAuthenticationException,
            TimeoutException {
        // share目录
        String shareDir = this.makeKey("share", env);
        EtcdResponsePromise<EtcdKeysResponse> resp1 = ec.getDir(shareDir).recursive().send();
        for (EtcdNode node : resp1.get().getNode().nodes) {
            entries.put(node.key.replaceAll(shareDir + "/", ""), node.value);
        }
        // 监听项目目录
        String watchDir = this.makeKey(prjName, env);
        EtcdResponsePromise<EtcdKeysResponse> resp2 = ec.getDir(watchDir).recursive().send();
        for (EtcdNode node : resp2.get().getNode().nodes) {
            entries.put(node.key.replaceAll(watchDir + "/", ""), node.value);
        }
    }

    public void dump() {
        String filename = prop.getProperty("dunpPath", "confs.dump");
        try (BufferedOutputStream out = new BufferedOutputStream(new FileOutputStream(filename))) {
            for (Entry<String, String> e : this.entries.entrySet()) {
                out.write((e.getKey() + "=" + e.getValue() + "\n").getBytes());
            }
            out.flush();
        } catch (Exception e) {
            LOG.error(e.getMessage());
        }
    }

    private void readFromDump() {
        try {
            String filename = prop.getProperty("dunpPath", "confs.dump");
            Properties p = new Properties();
            p.load(new FileInputStream(filename));
            for (Entry<Object, Object> e : p.entrySet()) {
                this.entries.put((String) e.getKey(), (String) e.getValue());
            }
        } catch (Exception e) {
            LOG.error(e.getMessage());
        }
    }

    private void watch(String path, Callback callback) {
        try {
            EtcdResponsePromise<EtcdKeysResponse> promise = ec.get(path).waitForChange().send();
            promise.addListener(p -> {
                callback.exec(p);
            });
        } catch (IOException e) {
            LOG.error(e.getMessage());
        }
    }

    public void watching(Callback callback) {
        watch(this.makeKey("publish", prjName, env), callback);
    }

    public void watchingShare(Callback callback) {
        watch(this.makeKey("share", env), callback);
    }

    public String get(String key) {
        return this.entries.get(key);
    }

    private void heart(String info) throws IOException {
        this.ec.put(this.makeKey("heartbeat", prjName, env, info), "1").ttl(5).send();
    }

    private String getInfo() {
        String server = prop.getProperty("instanceName", "instance") + "-";
        try {
            InetAddress addr = InetAddress.getLocalHost();
            server += addr.getHostName() + "-" + addr.getHostAddress();
        } catch (UnknownHostException e) {
            LOG.error(e.getMessage());
        }
        return server;
    }

    public String makeKey(String... paths) {
        String full = "";
        for (String path : paths) {
            full += "/" + path;
        }
        return full;
    }
}

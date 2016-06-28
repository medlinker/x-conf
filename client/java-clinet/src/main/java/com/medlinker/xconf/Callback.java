package com.medlinker.xconf;

import mousio.client.promises.ResponsePromise;
import mousio.etcd4j.responses.EtcdKeysResponse;


@FunctionalInterface
public interface Callback {
    public void exec(ResponsePromise<EtcdKeysResponse> p);
}

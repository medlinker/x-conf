package com.medlinker.xconf;

import org.hamcrest.Matchers;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

public class ClientTest {
    
    private Client c;
    
    
    @Before
    public void init() {
        c = new Client();
    }
    
    @Test
    public void TestGet() {
       Assert.assertThat(c.get("redis.host"), Matchers.notNullValue());
    }
}

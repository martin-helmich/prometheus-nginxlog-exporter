port = 4040

namespace "testone" {
  source_files = [".behave-sandbox/access.log"]
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
  labels {
    test1 = "1-1"
    test2 = "1-2"
  }
}

namespace "testtwo" {
  source_files = [".behave-sandbox/access-two.log"]
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
  labels {
    test1 = "2-1"
    test2 = "2-2"
  }
}


namespace "testthree" {
  source_files = [".behave-sandbox/access-three.log"]
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
  labels {
    test1 = "3-1"
    test2 = "3-2"
  }
}


namespace "testfour" {
  source_files = [".behave-sandbox/access-four.log"]
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
  labels {
    test1 = "4-1"
    test2 = "4-2"
  }
}


namespace "testfive" {
  source_files = [".behave-sandbox/access-five.log"]
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
  labels {
    test1 = "5-1"
    test2 = "5-2"
  }
}

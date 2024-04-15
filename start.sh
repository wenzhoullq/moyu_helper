export ENV="dev"
kill -s 9 `pgrep weixin_LLM`
chmod 777 weixin_LLM
nohup ./weixin_LLM &
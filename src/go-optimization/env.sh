python -m venv venv

source ./venv/scripts/activate
export PATH=`pwd`/venv/bin:$PATH

pip install mkdocs mkdocs-material mkdocs-git-revision-date-localized-plugin mkdocs-include-markdown-plugin mkdocs-rss-plugin

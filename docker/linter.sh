#!/bin/sh
case "$DS_BASE_VERSION" in
 @latest)  ;;
 ""    )  ;;
 *     ) echo "WARNING: DS_BASE_VERSION   is hardset to a value ($DS_BASE_VERSION). Maybe unstick it" ;;
esac

case "$DS_GITHUB_VERSION" in
  @latest)  ;;
 ""    )  ;;
 *     ) echo "WARNING: DS_GITHUB_VERSION is hardset to a value ($DS_GITHUB_VERSION). Maybe unstick it" ;;
esac

case "$DS_TUI_VERSION" in
  @latest)  ;;
 ""    )  ;;
 *     ) echo "WARNING: DS_TUI_VERSION    is hardset to a value ($DS_TUI_VERSION). Maybe unstick it" ;;
esac

case "$DS_GCLOUD_VERSION" in
  @latest)  ;;
 ""    )  ;;
 *     ) echo "WARNING: DS_GCLOUD_VERSION is hardset to a value ($DS_GCLOUD_VERSION). Maybe unstick it" ;;
esac

case "$DS_CONFIG_VERSION" in
  @latest)  ;;
 ""    )  ;;
 *     ) echo "WARNING: DS_CONFIG_VERSION is hardset to a value ($DS_CONFIG_VERSION). Maybe unstick it" ;;
esac
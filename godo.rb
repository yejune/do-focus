class Godo < Formula
  desc "CLI installer for Do - Claude Code project environment"
  homepage "https://github.com/yejune/do-focus"
  url "https://github.com/yejune/do-focus/archive/refs/tags/v0.1.57.tar.gz"
  sha256 "e884092216026f7a9a687897e494eb3fe42121dbb1e571dec902ec79d21220fd"
  license "MIT"
  head "https://github.com/yejune/do-focus.git", branch: "main"

  depends_on "go" => :build

  def install
        bin.install "godo"
    bin.install "do-worker"
  end

  def test
        assert_match "godo version", shell_output("#{bin}/godo version")
  end

  def caveats
    <<~EOS
            godo has been installed!
      
      Usage:
        cd my-project
        godo init      # Install Do
        godo update    # Update Do
    EOS
  end
end

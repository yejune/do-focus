class Godo < Formula
  desc "CLI installer for Do - Claude Code project environment"
  homepage "https://github.com/yejune/do-focus"
  url "https://github.com/yejune/do-focus/archive/refs/tags/v0.1.55.tar.gz"
  sha256 "2a1f8a7d72d4195a805c349d0a6b28efb61b19188a6fd49d3e25133da8b0bbdd"
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

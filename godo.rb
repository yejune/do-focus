class Godo < Formula
  desc "CLI installer for Do - Claude Code project environment"
  homepage "https://github.com/yejune/do-focus"
  url "https://github.com/yejune/do-focus/archive/refs/tags/v0.1.58.tar.gz"
  sha256 "43b389e9e3b93ba1f2bd6a4d7bae8d2b9076e3077ff362cbcfe0d889f4ee0b22"
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

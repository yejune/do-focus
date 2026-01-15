class Godo < Formula
  desc "CLI installer for Do - Claude Code project environment"
  homepage "https://github.com/yejune/do-focus"
  url "https://github.com/yejune/do-focus/archive/refs/tags/v0.1.69.tar.gz"
  sha256 "56cb1e3ba3b8b6058f5bed6b4662b203bb959bd4b4e9b6c6918d0dc9c37fdb4b"
  license "MIT"
  head "https://github.com/yejune/do-focus.git", branch: "main"

  depends_on "go" => :build

  def install
        system "go", "build", "-ldflags", "-X main.version=#{version}", "-o", "godo", "./cmd/godo/"
    Dir.chdir(".do/worker") do
      system "go", "build", "-ldflags", "-X main.version=#{version}", "-o", "../../godo-worker", "./cmd/worker/"
    end
    bin.install "godo"
    bin.install "godo-worker"
  end

  def test
        assert_match "godo version", shell_output("#{bin}/godo version")
  end

  def caveats
    <<~EOS
            godo has been installed!
      
      Usage:
        cd my-project
        godo sync             # Install or update Do
        godo worker start     # Start memory worker
    EOS
  end
end

class Godo < Formula
  desc "CLI installer for Do - Claude Code project environment"
  homepage "https://github.com/yejune/do-focus"
  url "https://github.com/yejune/do-focus/archive/refs/tags/v0.1.61.tar.gz"
  sha256 "867b03d3e407dc6969b8fe2e199d44ea4ea70f88908741d309012f5b281a5291"
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
